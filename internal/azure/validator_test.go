package azure_test

import (
	"encoding/json"
	"github.com/ms-henglu/terraform-provider-azurermg/utils"
	"testing"

	"github.com/ms-henglu/terraform-provider-azurermg/internal/azure"
)

func Test_BodyValidation(t *testing.T) {
	testData := []struct {
		Id         string
		ApiVersion string
		Body       string
		Error      bool
	}{
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

	for _, data := range testData {
		resourceType := utils.GetResourceType(data.Id)

		var body interface{}
		json.Unmarshal([]byte(data.Body), &body)

		def, err := azure.GetResourceDefinition(resourceType, data.ApiVersion)
		if err != nil {
			t.Fatal(err)
		}

		if def != nil {
			errors := (*def).Validate(body, "")
			if (len(errors) > 0) != data.Error {
				t.Errorf("expect error: %t, got error: %t for id: %s, api-version: %s", data.Error, len(errors) > 0, data.Id, data.ApiVersion)
			}
		} else {
			t.Fatalf("failed to load resource definition for id: %s, api-version: %s", data.Id, data.ApiVersion)
		}
	}
}
