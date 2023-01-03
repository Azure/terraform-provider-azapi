package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ActionResource struct{}

func TestAccActionResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := ActionResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func (r ActionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azapi_resource_action" "list" {
  type                   = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id            = azapi_resource.test.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id = azapi_resource.test.id
  action      = "agentRegistrationInformation/regenerateKey"
  body = jsonencode({
    keyName = "primary"
  })
  depends_on = [
    data.azapi_resource_action.list
  ]
}
`, GenericResource{}.defaultTag(data))
}
