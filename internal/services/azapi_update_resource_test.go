package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type GenericUpdateResource struct{}

func TestAccGenericUpdateResource_automationAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.automationAccount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccGenericUpdateResource_withNameParentId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.automationAccountWithNameParentId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccGenericUpdateResource_siteConfigSlotConfigNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.siteConfigSlotConfigNames(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccGenericUpdateResource_locks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.locks(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccGenericUpdateResource_ignoreOrderInArray(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreOrderInArray(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericUpdateResource_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.timeouts(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericUpdateResource_ignoreIDCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreIDCasing(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.ignoreIDCasingUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericUpdateResource_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.headers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericUpdateResource_queryParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.queryParameters(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r GenericUpdateResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resourceType := state.Attributes["type"]
	id, err := parse.ResourceIDWithResourceType(state.ID, resourceType)
	if err != nil {
		return nil, err
	}

	resp, err := client.ResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion, clients.DefaultRequestOptions())
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			exist := false
			return &exist, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	exist := utils.GetId(resp) != nil
	return &exist, nil
}

func (r GenericUpdateResource) automationAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azapi_resource.automationAccount.id
  body = {
    properties = {
      publicNetworkAccess = true
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) automationAccountWithNameParentId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_update_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = azapi_resource.automationAccount.name
  parent_id = azapi_resource.resourceGroup.id
  body = {
    properties = {
      publicNetworkAccess = true
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) siteConfigSlotConfigNames(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s


resource "azapi_resource" "serverfarm" {
  type      = "Microsoft.Web/serverfarms@2023-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      hyperV         = false
      perSiteScaling = false
      reserved       = false
      zoneRedundant  = false
    }
    sku = {
      name = "S1"
    }
  }
}

resource "azapi_resource" "site" {
  type      = "Microsoft.Web/sites@2023-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      clientAffinityEnabled = false
      clientCertEnabled     = false
      clientCertMode        = "Required"
      enabled               = true
      httpsOnly             = false
      publicNetworkAccess   = "Enabled"
      serverFarmId          = azapi_resource.serverfarm.id
      siteConfig = {
        acrUseManagedIdentityCreds       = false
        alwaysOn                         = true
        autoHealEnabled                  = false
        ftpsState                        = "Disabled"
        http20Enabled                    = false
        loadBalancing                    = "LeastRequests"
        localMySqlEnabled                = false
        managedPipelineMode              = "Integrated"
        minTlsVersion                    = "1.2"
        publicNetworkAccess              = "Enabled"
        remoteDebuggingEnabled           = false
        scmIpSecurityRestrictionsUseMain = false
        scmMinTlsVersion                 = "1.2"
        use32BitWorkerProcess            = true
        vnetRouteAllEnabled              = false
        webSocketsEnabled                = false
        windowsFxVersion                 = ""
      }
      vnetRouteAllEnabled = false
    }
  }
}

resource "azapi_update_resource" "test" {
  type      = "Microsoft.Web/sites/config@2021-03-01"
  name      = "slotConfigNames"
  parent_id = azapi_resource.site.id
  body = {
    properties = {
      connectionStringNames   = ["test1", "test2"]
      appSettingNames         = ["test3"]
      azureStorageConfigNames = ["test4"]
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) locks(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_update_resource" "test1" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azapi_resource.automationAccount.id
  body = {
    properties = {
      publicNetworkAccess = true
    }
  }
  locks = [azapi_resource.automationAccount.id]
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azapi_resource.automationAccount.id
  body = {
    properties = {
      publicNetworkAccess = true
    }
  }
  locks = [azapi_resource.automationAccount.id]
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) ignoreOrderInArray(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "vnet" {
  type      = "Microsoft.Network/virtualNetworks@2023-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/16",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
        {
          name = "first"
          properties = {
            addressPrefix = "10.0.3.0/24"
          }
        },
        {
          name = "second"
          properties = {
            addressPrefix = "10.0.4.0/24"
          }
        }
      ]
    }
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Network/virtualNetworks@2022-07-01"
  resource_id = azapi_resource.vnet.id
  body = {
    properties = {
      subnets = [
        {
          name = "second"
          properties = {
            addressPrefix = "10.0.4.0/24"
          }
        },
        {
          name = "first"
          properties = {
            addressPrefix = "10.0.3.0/24"
          }
        }
      ]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r GenericUpdateResource) ignoreIDCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_update_resource" "test" {
  type        = "Microsoft.SignalRService/WebPubSub@2024-04-01-preview"
  resource_id = azapi_resource.webPubSub.id

  body = {
    properties = {
      networkACLs = {
        defaultAction = "Deny"
        publicNetwork = {
          allow = ["ClientConnection"]
        }
        ipRules = [{
          value  = "0.0.0.0/0"
          action = "Allow"
        }]
      }
    }
  }
}
`, r.templateForIDCasing(data))
}

func (r GenericUpdateResource) ignoreIDCasingUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_update_resource" "test" {
  type        = "Microsoft.SignalRService/WebPubSub@2024-04-01-preview"
  resource_id = azapi_resource.webPubSub.id

  body = {
    properties = {
      networkACLs = {
        defaultAction = "Deny"
        publicNetwork = {
          allow = ["ClientConnection", "RESTAPI"]
        }
        ipRules = [{
          value  = "0.0.0.0/0"
          action = "Allow"
        }]
      }
    }
  }
}
`, r.templateForIDCasing(data))
}

func (GenericUpdateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString)
}

func (r GenericUpdateResource) templateForIDCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "webPubSub" {
  type      = "Microsoft.SignalRService/webPubSub@2024-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      disableAadAuth      = false
      disableLocalAuth    = false
      publicNetworkAccess = "Enabled"
      tls = {
        clientCertEnabled = false
      }
    }
    sku = {
      capacity = 1
      name     = "Standard_S1"
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) timeouts(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azapi_resource.automationAccount.id
  body = {
    properties = {
      publicNetworkAccess = true
    }
  }
  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
    read   = "10m"
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) oldConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azapi_resource.automationAccount.id
  body = jsonencode({
    properties = {
      publicNetworkAccess = true
    }
  })
  response_export_values = ["properties"]
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) headers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azapi_resource.automationAccount.id
  body = {
    properties = {
      publicNetworkAccess = true
    }
  }
  update_headers = {
    "header2" = "update-value"
  }
  read_headers = {
    "header4" = "read-value"
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) queryParameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azapi_resource.automationAccount.id
  body = {
    properties = {
      publicNetworkAccess = true
    }
  }
  update_query_parameters = {
    "query1" = ["update-value"]
  }
  read_query_parameters = {
    "query1" = ["read-value"]
  }
}
`, r.template(data), data.RandomString)
}
