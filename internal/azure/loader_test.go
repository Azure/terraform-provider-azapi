package azure_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure"
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
