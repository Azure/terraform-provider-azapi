package azure_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
)

func Test_GetAzureSchema(t *testing.T) {
	if azure.GetAzureSchema() == nil {
		t.Errorf("failed to load azure schema")
	}
}

func Test_GetApiVersions(t *testing.T) {
	case1 := "Microsoft.MachineLearningServices/workspaces/computes"
	if len(azure.GetApiVersions(case1)) == 0 {
		t.Errorf("expect multiple api-version but got 0 for Microsoft.MachineLearningServices/workspaces/computes")
	}

	case2 := "Microsoft.MachineLearningServices/workspaces/computes0"
	if len(azure.GetApiVersions(case2)) != 0 {
		t.Errorf("expect 0 api-version but got multiple for Microsoft.MachineLearningServices/workspaces/computes0")
	}
}

func Test_GetResourceDefinition(t *testing.T) {
	case1 := "Microsoft.MachineLearningServices/workspaces/computes"
	versions := azure.GetApiVersions(case1)
	for _, v := range versions {
		def, err := azure.GetResourceDefinition(case1, v)
		if err != nil {
			t.Error(err)
		}
		if def == nil {
			t.Errorf("failed to load resource definition for %s api-version %s", case1, v)
		}
	}
}

func Test_AllBicepTypes(t *testing.T) {
	if schema := azure.GetAzureSchema(); schema == nil {
		t.Fatal("failed to load azure schema")
	} else {
		resources := schema.Resources
		if len(resources) == 0 {
			t.Fatal("expect resources are not empty")
		}
		for resourceName, res := range resources {
			if len(resourceName) == 0 {
				t.Fatal("expect resource name is not empty")
			}
			if res == nil {
				t.Fatalf("expect resource definition is not nil, resource name: %s", resourceName)
			} else {
				definitions := res.Definitions
				if len(definitions) == 0 {
					t.Fatalf("expect resource definitions are not empty, resource name: %s", resourceName)
				}
				for _, definition := range definitions {
					_, err := definition.GetDefinition()
					if err != nil {
						t.Fatal(err)
					}
				}
			}
		}
	}
}
