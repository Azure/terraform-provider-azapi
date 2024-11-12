package services_test

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExamples(t *testing.T) {
	cases := strings.Split(os.Getenv("ARM_TEST_EXAMPLES"), ",")

	for _, c := range cases {
		if strings.Trim(c, " ") == "" {
			continue
		}
		t.Run(c, func(t *testing.T) {
			workingDir := path.Join("..", "..", "examples", c)
			t.Logf("Running example %s", workingDir)

			content, err := os.ReadFile(path.Join(workingDir, "main.tf"))
			if err != nil {
				t.Errorf("Error reading main.tf: %v", err)
				return
			}

			data := acceptance.BuildTestData(t, "azapi_resource", "test")
			config := string(content)
			config = strings.ReplaceAll(config, `terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}`, "")
			config = strings.ReplaceAll(config, `default = "acctest0001"`, fmt.Sprintf(`default = "acctest%s"`, data.RandomString))

			r := GenericResource{}
			data.ResourceTest(t, r, []resource.TestStep{
				{
					Config: config,
				},
			})
		})
	}
}
