package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type GenericOperationResource struct{}

func TestAccGenericOperationResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_operation", "test")
	r := GenericOperationResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func (r GenericOperationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s


data "azapi_operation" "list" {
  type                   = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id            = azapi_resource.test.id
  operation              = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_operation" "test" {
  type        = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id = azapi_resource.test.id
  operation   = "agentRegistrationInformation/regenerateKey"
  body = jsonencode({
    keyName = "primary"
  })
  depends_on = [
    data.azapi_operation.list
  ]
}
`, GenericResource{}.defaultTag(data))
}
