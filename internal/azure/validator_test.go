package azure_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ms-henglu/terraform-provider-azurermg/internal/azure"
	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

func Test_BodyValidation(t *testing.T) {
	testData := []struct {
		Id         string
		ApiVersion string
		Body       string
		Error      bool
	}{{
		Id:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/compute1",
		ApiVersion: "2021-07-01",
		Body: `
{
    "location": "eastus",
    "properties": {
        "properties": {
            "state": "Running"
        }
    }
}
`,
		Error: true, // properties.computeType is required
	},
		{
			Id:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/compute1",
			ApiVersion: "2021-07-01",
			Body: `
{
    "location": "eastus",
    "properties": {
        "computeType": "ComputeInstance1",
        "properties": {
            "state": "Running"
        }
    }
}
`,
			Error: true, // invalid properties.computeType
		},
		{
			Id:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/compute1",
			ApiVersion: "2021-07-01",
			Body: `
{
    "location": "eastus",
    "properties": {
        "computeType": "ComputeInstance",
        "properties": {
            "state": "Running"
        }
    }
}
`,
			Error: true, // properties.properties.state is read only
		},
		{
			Id:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/compute1",
			ApiVersion: "2021-07-01",
			Body: `
{
    "location": "eastus",
    "properties": {
        "computeType": "ComputeInstance",
        "properties": {
        }
    }
}
`,
			Error: false,
		},
		{
			Id:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-211109150453866525/providers/Microsoft.ContainerRegistry/registries/acctest61311",
			ApiVersion: "2020-11-01-preview",
			Body: `
{
  "identity": {
    "type": "SystemAssigned, UserAssigned",
    "userAssignedIdentities": {
      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-211109152718418172/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctestb0i47": {}
    }
  },
  "location": "westeurope",
  "properties": {
    "adminUserEnabled": true
  },
  "sku": {
    "name": "Standard"
  },
  "tags": {
    "Key": "Value"
  }
}
`,
			Error: false,
		},
		{
			Id:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.AppPlatform/Spring/myservice/apps/myapp",
			ApiVersion: "2020-07-01",
			Body: `
{
    "location": "westeurope",
    "properties": {
        "activeDeploymentName": "mydeployment1",
        "fqdn": "myapp.mydomain.com",
        "httpsOnly": 1,
        "persistentDisk": {
            "mountPath": "/persistent",
            "sizeInGB": 2
        },
        "publi1c": true,
        "temporaryDisk": {
            "mountPath": "/temp",
            "sizeInGB": 2
        }
    }
}
`,
			Error: true,
		},
		{
			Id:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Automation/automationAccounts/acctest3093",
			ApiVersion: "2021-06-22",
			Body: `
{
    "location": "westeurope",
    "name": "myAutomationAccount9",
    "properties": {
        "sku": {
            "name": "Free"
        }
    }
}
`,
			Error: false,
		},
	}

	for index, data := range testData {
		resourceType := utils.GetResourceType(data.Id)

		var body interface{}
		_ = json.Unmarshal([]byte(data.Body), &body)

		def, err := azure.GetResourceDefinition(resourceType, data.ApiVersion)
		if err != nil {
			t.Fatal(err)
		}

		if def != nil {
			errors := (*def).Validate(body, "")
			fmt.Printf("Running test for case %d, resource type: %s, api-version: %s\n", index, resourceType, data.ApiVersion)
			fmt.Println(errors)
			if (len(errors) > 0) != data.Error {
				t.Errorf("expect error: %t, got error: %t for id: %s, api-version: %s", data.Error, len(errors) > 0, data.Id, data.ApiVersion)
			}
		} else {
			t.Fatalf("failed to load resource definition for id: %s, api-version: %s", data.Id, data.ApiVersion)
		}
	}
}
