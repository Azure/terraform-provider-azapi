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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccGenericUpdateResource_listUniqueIdProperty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:            r.listUniqueIdProperty(data),
			ExternalProviders: externalProvidersAzurerm(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestCheckOutput("azure_policy_evaluation_details_enabled", "true"),
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

func TestAccGenericUpdateResource_SensitiveBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.SensitiveBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericUpdateResource_SensitiveBodyVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.SensitiveBodyWithHash(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.tags.tag1").HasValue("tag1-value"),
			),
		},
		{
			Config: r.SensitiveBodyWithVersion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.tags.tag1").HasValue("tag1-value"),
				check.That(data.ResourceName).Key("output.tags.tag2").HasValue("tag2-value2"),
			),
		},
		{
			Config: r.SensitiveBodyWithVersionMultipleTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.tags.tag1").HasValue("tag1-value"),
				check.That(data.ResourceName).Key("output.tags.tag2").HasValue("tag2-value2"),
				check.That(data.ResourceName).Key("output.tags.tag3").DoesNotExist(),
			),
		},
		{
			Config: r.SensitiveBodyWithHashMultipleTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.tags.tag1").HasValue("tag1-value"),
				check.That(data.ResourceName).Key("output.tags.tag2").HasValue("tag2-value3"),
				check.That(data.ResourceName).Key("output.tags.tag3").HasValue("tag3-value"),
			),
		},
	})
}

func TestAccGenericUpdateResource_sensitiveBodyVersionWithEmptyBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.sensitiveBodyVersionWithEmptyBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericUpdateResource_BadUserAssignedIdentitiesSchema(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.BadUserAssignedIdentitiesSchema(data),
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


resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
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
      ]
    }
  }
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "acctest%[2]d"
  body = {
    properties = {
      addressPrefix = "10.0.2.0/24"
      delegations = [
      ]
      defaultOutboundAccess             = false
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
      ]
    }
  }
}

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          name = "testconfiguration1"
          properties = {
            primary                   = true
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        },
        {
          name = "testconfiguration2"
          properties = {
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        },
      ]
    }
  }

  lifecycle {
    ignore_changes = [
      body.properties.ipConfigurations[0].properties.primary,
      body.properties.ipConfigurations[1].properties.primary,
    ]
  }

}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Network/networkInterfaces@2022-07-01"
  resource_id = azapi_resource.networkInterface.id
  body = {
    properties = {
      ipConfigurations = [
        {
          name = "testconfiguration2"
          properties = {
            primary = true
          }
        },
        {
          name = "testconfiguration1"
          properties = {
            primary = false
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

func (r GenericUpdateResource) SensitiveBody(data acceptance.TestData) string {
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
  sensitive_body = {
    properties = {
      publicNetworkAccess = true
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) SensitiveBodyWithHash(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
    }
  }
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.DataFactory/factories@2018-06-01"
  resource_id = azapi_resource.factory.id
  body = {
    tags = {
      tag1 = "tag1-value"
    }
  }
  sensitive_body = {
    tags = {
      tag2 = "tag2-value"
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) SensitiveBodyWithVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
    }
  }
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.DataFactory/factories@2018-06-01"
  resource_id = azapi_resource.factory.id
  body = {
    tags = {
      tag1 = "tag1-value"
    }
  }
  sensitive_body = {
    tags = {
      tag2 = "tag2-value2"
    }
  }

  sensitive_body_version = {
    "tags.tag2" = "2"
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) SensitiveBodyWithVersionMultipleTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
    }
  }
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.DataFactory/factories@2018-06-01"
  resource_id = azapi_resource.factory.id
  body = {
    tags = {
      tag1 = "tag1-value"
    }
  }
  sensitive_body = {
    tags = {
      tag2 = "tag2-value3"
      tag3 = "tag3-value"
    }
  }

  sensitive_body_version = {
    "tags.tag2" = "2"
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) SensitiveBodyWithHashMultipleTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "factory" {
  type      = "Microsoft.DataFactory/factories@2018-06-01"
  name      = "acctest-%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      repoConfiguration   = null
    }
  }
  lifecycle {
    ignore_changes = [tags]
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.DataFactory/factories@2018-06-01"
  resource_id = azapi_resource.factory.id
  body = {
    tags = {
      tag1 = "tag1-value"
    }
  }
  sensitive_body = {
    tags = {
      tag2 = "tag2-value3"
      tag3 = "tag3-value"
    }
  }
}
`, r.template(data), data.RandomString)
}

// sensitiveBodyVersionWithEmptyBody tests issue #999 scenario
func (r GenericUpdateResource) sensitiveBodyVersionWithEmptyBody(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2024-10-23"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-%[2]s"
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
  type      = "Microsoft.Automation/automationAccounts@2024-10-23"
  parent_id = azapi_resource.resourceGroup.id
  name      = azapi_resource.automationAccount.name

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }

  sensitive_body = {

  }

  sensitive_body_version = {
    "properties.publicNetworkAccess" = "1"
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) BadUserAssignedIdentitiesSchema(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s


resource "azapi_resource" "managedIdentity1" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30"
  parent_id = azapi_resource.resourceGroup.id
  name      = "actest-%[2]s"
  location  = azapi_resource.resourceGroup.location
  body      = {}
}

resource "azapi_resource" "managedIdentity2" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30"
  parent_id = azapi_resource.resourceGroup.id
  name      = "actest2-%[2]s"
  location  = azapi_resource.resourceGroup.location
  body      = {}
}


resource "azapi_resource" "apiManagementInstance" {
  type      = "Microsoft.ApiManagement/service@2020-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-%[2]s"
  location  = azapi_resource.resourceGroup.location
  identity {
    type = "UserAssigned"
    identity_ids = [
      azapi_resource.managedIdentity1.id,
      azapi_resource.managedIdentity2.id
    ]
  }

  body = {
    sku = {
      capacity = 1
      name     = "Developer"
    }
    properties = {
      virtualNetworkType = "None"
      publisherEmail     = "publisherEmail@contoso.com"
      publisherName      = "publisherName"
    }
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.ApiManagement/service@2020-12-01"
  resource_id = azapi_resource.apiManagementInstance.id
  body = {
    properties = {
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) listUniqueIdPropertyTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azapi_client_config" "current" {
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2025-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestvault%[3]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      enableRbacAuthorization   = true
      enableSoftDelete          = true
      publicNetworkAccess       = "Enabled"
      softDeleteRetentionInDays = 7
      sku = {
        family = "A"
        name   = "standard"
      }      
      tenantId = data.azapi_client_config.current.tenant_id
    }
  }
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2025-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestlaw%[3]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku                        = { name = "PerGB2018" }
      retentionInDays            = 30
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
    }
  }
}

# Create the diagnostic setting with basic config first
resource "azapi_resource" "diagnosticSetting" {
  type      = "Microsoft.Insights/diagnosticSettings@2021-05-01-preview"
  parent_id = azapi_resource.vault.id
  name      = "acctest%[2]d"
  body = {
    properties = {
      workspaceId = azapi_resource.workspace.id
      logs = [
        {
          category = "AuditEvent"
          enabled  = true
        },
        {
          category = "AzurePolicyEvaluationDetails"
          enabled  = false
        }
      ]
    }
  }

  ignore_missing_property = true

  # Ignore changes to logs since azapi_update_resource will manage them
  lifecycle {
    ignore_changes = [body.properties.logs]
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r GenericUpdateResource) listUniqueIdProperty(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

# Use azapi_update_resource to update specific log categories
# without affecting others that Azure may have added
resource "azapi_update_resource" "test" {
  type        = "Microsoft.Insights/diagnosticSettings@2021-05-01-preview"
  resource_id = azapi_resource.diagnosticSetting.id
  body = {
    properties = {
      logs = [
        {
          category = "AzurePolicyEvaluationDetails"
          enabled  = true
        }
      ]
    }
  }

  # Use composite key to match log entries by both category and categoryGroup
  # This handles cases where Azure uses either field to identify a log setting
  list_unique_id_property = {
    "properties.logs" = "category, categoryGroup"
  }

  # Only update the logs we specify, ignore any others
  ignore_other_items_in_list = ["properties.logs"]

  response_export_values = ["properties.logs"]
}

locals {
  logs = azapi_update_resource.test.output.properties.logs
  azure_policy_evaluation_details_enabled = try([for l in local.logs : l.enabled if l.category == "AzurePolicyEvaluationDetails"][0], null)
}

output "azure_policy_evaluation_details_enabled" {
  value = tostring(local.azure_policy_evaluation_details_enabled)
}
`, r.listUniqueIdPropertyTemplate(data), data.RandomInteger)
}
