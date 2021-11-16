package services_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/acceptance"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/acceptance/check"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/clients"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/parse"
	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

type GenericPatchResource struct{}

func TestAccGenericPatchResource_loadBalancerNatRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurermg_patch_resource", "test")
	r := GenericPatchResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.loadBalancerNatRule(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericPatchResource_sqlServer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurermg_patch_resource", "test")
	r := GenericPatchResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.sqlServer(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericPatchResource_springCloudApp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurermg_patch_resource", "test")
	r := GenericPatchResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.springCloudApp(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r GenericPatchResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, response, err := client.ResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if response.StatusCode == http.StatusNotFound {
			exist := false
			return &exist, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	exist := len(utils.GetId(resp)) != 0
	return &exist, nil
}

func (r GenericPatchResource) loadBalancerNatRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_public_ip" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_nat_rule" "test" {
  resource_group_name            = azurerm_resource_group.test.name
  loadbalancer_id                = azurerm_lb.test.id
  name                           = "RDPAccess"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "PublicIPAddress"
}

resource "azurermg_patch_resource" "test" {
  resource_id = azurerm_lb.test.id
  type        = "Microsoft.Network/loadBalancers@2021-03-01"
  body        = <<BODY
    {
      "properties": {
        "inboundNatRules": [
          {
            "properties": {
               "idleTimeoutInMinutes": 15
            }
          }
        ]
      }
    }
    BODY

  depends_on = [
    azurerm_lb_nat_rule.test,
  ]
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (r GenericPatchResource) sqlServer(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_sql_server" "test" {
  name                         = "acctest-%[2]s"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurermg_patch_resource" "test" {
  resource_id = azurerm_sql_server.test.id
  type        = "Microsoft.Sql/servers@2021-02-01-preview"
  body        = <<BODY
{
  "properties": {
    "restrictOutboundNetworkAccess": "Disabled"
  }
}
  BODY
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (r GenericPatchResource) springCloudApp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  service_name        = azurerm_spring_cloud_service.test.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurermg_patch_resource" "test" {
  resource_id = azurerm_spring_cloud_app.test.id
  type        = "Microsoft.AppPlatform/Spring/apps@2021-06-01-preview"
  body        = <<BODY
  {
    "properties": {
      "temporaryDisk": {
        "mountPath": "/temp",
        "sizeInGB": 4
      }
    }
  }
  BODY
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (GenericPatchResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.LocationPrimary, data.RandomStringOfLength(10))
}
