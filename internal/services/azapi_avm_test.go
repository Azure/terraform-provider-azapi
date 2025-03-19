package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/provider"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zclconf/go-cty/cty"
)

func TestAccAzapiRelease_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_AVM") == "" {
		t.Skip("Skipping AVM tests as ARM_TEST_AVM is not set")
	}
	modulesName := []string{
		"terraform-azurerm-avm-res-network-virtualnetwork",
	}
	skipped := map[string]map[string]bool{}

	wd, _ := os.Getwd()
	for _, moduleName := range modulesName {
		t.Logf("Running regression test for %s", moduleName)

		repoDirectory := path.Join(wd, "modules", moduleName)
		err := CloneModule(moduleName, repoDirectory)
		if err != nil {
			t.Fatalf("Error cloning AVM module: %v", err)
		}

		testcases, err := ListExamples(repoDirectory)
		if err != nil {
			t.Fatalf("Error listing examples: %v", err)
		}

		for _, testcase := range testcases {
			if skipped[moduleName] != nil && skipped[moduleName][testcase] {
				t.Logf("Skipping test for %s", testcase)
				continue
			}
			t.Logf("Running test for %s", testcase)
			ModuleUpgradeTest(t, repoDirectory, testcase)
		}
	}
}

func ModuleUpgradeTest(t *testing.T, repoDirectory string, caseName string) {
	configDir := path.Join(repoDirectory, "examples", caseName)
	workingDirectory := path.Join(repoDirectory, "azapi")

	// Create working directory for the test
	err := os.MkdirAll(workingDirectory, 0755)
	if err != nil {
		t.Fatalf("Error creating working directory: %v", err)
	}

	// Update the azurerm_resource_group name to have an "acctest" prefix, to bypass the subscription's policy checks
	testdata := acceptance.BuildTestData(t, "azapi_resource", "test")
	resourceGroupName := "acctest" + testdata.RandomStringOfLength(5)
	err = UpdateConfigDirectory(configDir, map[string]map[string]cty.Value{
		"resource.azurerm_resource_group": {
			"name": cty.StringVal(resourceGroupName),
		},
	})
	if err != nil {
		t.Fatalf("Error updating config directory: %v", err)
	}

	// Disable default output for azapi provider
	_ = os.Setenv("ARM_DISABLE_DEFAULT_OUTPUT", "true")

	defer func() {
		// testing-framework doesn't support cleanup resources if there are modules in the config
		// we need to clean up the resource group manually if error occurs
		client, err := acceptance.BuildTestClient()
		if err != nil {
			t.Errorf("Error building test client: %v", err)
		}

		resourceGroupId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", client.Account.GetSubscriptionId(), resourceGroupName)
		_, _ = client.ResourceClient.Delete(context.Background(), resourceGroupId, "2020-06-01", clients.DefaultRequestOptions())
	}()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				// step1: Deploy the config with latest released azapi
				ConfigDirectory: func(request config.TestStepConfigRequest) string {
					return configDir
				},
				ExternalProviders: map[string]resource.ExternalProvider{
					// it's not allowed to set external providers in the test step if config directory is set
					// by default the latest azapi provider will be used
				},
			},
			{
				// step2: Use the developing azapi provider to run `terraform plan` command and check if the plan is empty
				PlanOnly: true,
				ConfigDirectory: func(request config.TestStepConfigRequest) string {
					return configDir
				},
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"azapi": providerserver.NewProtocol6WithError(provider.AzureProvider()),
				},
			},
			{
				// step3: Destroy the resources. Because the framework-testing doesn't support destroying resources with modules.
				Destroy: true,
				ConfigDirectory: func(request config.TestStepConfigRequest) string {
					return configDir
				},
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"azapi": providerserver.NewProtocol6WithError(provider.AzureProvider()),
				},
			},
		},
		WorkingDir: workingDirectory,
	})
}

func UpdateConfigDirectory(dir string, updateValueMap map[string]map[string]cty.Value) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".tf") {
			continue
		}

		data, err := os.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			return err
		}

		hclFile, diags := hclwrite.ParseConfig(data, file.Name(), hcl.InitialPos)
		if diags.HasErrors() {
			return diags
		}

		for _, block := range hclFile.Body().Blocks() {
			key := block.Type()
			if len(block.Labels()) != 0 {
				key += "." + block.Labels()[0]
			}
			updateMap := updateValueMap[key]
			if updateMap == nil {
				continue
			}

			for k, v := range updateMap {
				block.Body().SetAttributeValue(k, v)
			}
		}

		err = os.WriteFile(path.Join(dir, file.Name()), hclFile.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func CloneModule(moduleName, repoDirectory string) error {
	_, err := git.PlainClone(repoDirectory, false, &git.CloneOptions{
		URL:      fmt.Sprintf("https://github.com/Azure/%s.git", moduleName),
		Progress: os.Stdout,
	})
	if err != nil && !errors.Is(err, git.ErrRepositoryAlreadyExists) {
		return err
	}
	return nil
}

func ListExamples(repoDirectory string) ([]string, error) {
	examplesDir := path.Join(repoDirectory, "examples")
	files, err := os.ReadDir(examplesDir)
	if err != nil {
		return nil, err
	}
	examples := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			examples = append(examples, file.Name())
		}
	}
	return examples, nil
}
