package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ActionDataSource struct{}

func TestAccActionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func (r ActionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_action" "test" {
  type                   = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id            = azapi_resource.test.id
  action                 = "listKeys"
  response_export_values = ["*"]
}
`, GenericResource{}.defaultTag(data))
}
