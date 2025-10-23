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

type DataPlaneResource struct{}

func TestAccDataPlaneResource_appConfigKeyValues(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.appConfigKeyValues(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_purviewClassification(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.purviewClassification(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_purviewCollection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.purviewCollection(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_keyVaultIssuer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.keyVaultIssuer(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_keyVaultSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.keyVaultSecret(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataPlaneResource{}),
			),
		},
		{
			Config: r.keyVaultSecretUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataPlaneResource{}),
			),
		},
	})
}

func TestAccDataPlaneResource_iotAppsUser(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.iotAppsUser(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.timeouts(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_searchServiceIndex(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.searchServiceIndex(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_searchServiceDataSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.searchServiceDataSource(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_searchServiceIndexer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.searchServiceIndexer(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_searchServiceSynonymMap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.searchServiceSynonymMap(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.headers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_queryParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.queryParameters(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_replaceTriggeredByExternalValues(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.replaceTriggeredByExternalValues(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (DataPlaneResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resourceType := state.Attributes["type"]
	id, err := parse.DataPlaneResourceIDWithResourceType(state.ID, resourceType)
	if err != nil {
		return nil, err
	}

	_, err = client.DataPlaneClient.Get(ctx, id, clients.DefaultRequestOptions())
	if err == nil {
		b := true
		return &b, nil
	}
	if utils.ResponseErrorWasNotFound(err) {
		b := false
		return &b, nil
	}
	return nil, fmt.Errorf("checking for presence of existing %s: %+v", id, err)
}

func (r DataPlaneResource) appConfigKeyValues(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "appconf" {
  type      = "Microsoft.AppConfiguration/configurationStores@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    sku = {
      name = "standard"
    }
  }
  response_export_values = {
    endpoint = "properties.endpoint"
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    appConfigDataOwnerRoleId = "value[?properties.roleName == 'App Configuration Data Owner'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.appconf.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.appConfigDataOwnerRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.AppConfiguration/configurationStores/keyValues@1.0"
  parent_id = replace(azapi_resource.appconf.output.endpoint, "https://", "")
  name      = "mykey"
  body = {
    content_type = ""
    value        = "myvalue"
  }


  depends_on = [
    azapi_resource.roleAssignment,
  ]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) purviewClassification(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
    }
  }
  response_export_values = ["*"]
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Purview/accounts/Scanning/classificationrules@2022-07-01-preview"
  parent_id = replace(azapi_resource.account.output.properties.endpoints.scan, "https://", "")
  name      = "acctest%[2]s"
  body = {
    kind = "Custom"
    properties = {
      description        = "Let's put a cool desc here"
      classificationName = "MICROSOFT.FINANCIAL.AUSTRALIA.BANK_ACCOUNT_NUMBER"
      columnPatterns = [
        {
          pattern = "^data$"
          kind    = "Regex"
        }
      ]
      dataPatterns = [
        {
          pattern = "^[0-9]{2}-[0-9]{4}-[0-9]{6}-[0-9]{3}$"
          kind    = "Regex"
        }
      ]
      minimumPercentageMatch = 60
      ruleStatus             = "Enabled"
    }
  }
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) purviewCollection(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
    }
  }
  response_export_values = ["*"]
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Purview/accounts/Account/collections@2019-11-01-preview"
  parent_id = "${azapi_resource.account.name}.purview.azure.com"
  name      = "defaultResourceSetRuleConfig"
  body = {
    friendlyName = "Finance"
  }
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) keyVaultIssuer(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azapi_client_config" "current" {}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        family = "A"
        name   = "standard"
      }
      tenantId                  = data.azapi_client_config.current.tenant_id
      enabledForDiskEncryption  = true
      softDeleteRetentionInDays = 7
      enablePurgeProtection     = true
      accessPolicies            = []
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [body.properties.accessPolicies]
  }
  response_export_values = {
    vaultUri = "properties.vaultUri"
  }
}

resource "azapi_resource_action" "add_accesspolicy" {
  type        = "Microsoft.KeyVault/vaults/accessPolicies@2023-02-01"
  resource_id = "${azapi_resource.vault.id}/accessPolicies/add"
  method      = "PUT"
  body = {
    properties = {
      accessPolicies = [{
        tenantId = data.azapi_client_config.current.tenant_id
        objectId = data.azapi_client_config.current.object_id
        permissions = {
          certificates = ["managecontacts", "getissuers", "setissuers", "deleteissuers"]
        }
      }]
    }
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.KeyVault/vaults/certificates/issuers@7.4"
  parent_id = trimsuffix(trimprefix(azapi_resource.vault.output.vaultUri, "https://"), "/")
  name      = "acctest%[2]s"
  body = {
    provider = "Test"
    credentials = {
      account_id = "keyvaultuser"
    }
    org_details = {
      admin_details = [
        {
          first_name = "John"
          last_name  = "Doe"
          email      = "admin@microsoft.com"
          phone      = "4255555555"
        }
      ]
    }
  }

  depends_on = [
    azapi_resource_action.add_accesspolicy
  ]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) keyVaultSecret(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.KeyVault/vaults/secrets@7.4"
  parent_id = trimsuffix(trimprefix(azapi_resource.vault.output.vaultUri, "https://"), "/")
  name      = "acctest%[2]s"
  body = {
    value = "secret-value"
  }

  depends_on = [
    azapi_resource_action.add_accesspolicy_secret
  ]
}`, r.keyVaultSecretTemplate(data), data.RandomString)
}

func (r DataPlaneResource) keyVaultSecretUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.KeyVault/vaults/secrets@7.4"
  parent_id = trimsuffix(trimprefix(azapi_resource.vault.output.vaultUri, "https://"), "/")
  name      = "acctest%[2]s"
  body = {
    value = "updated-secret-value"
    attributes = {
      enabled = true
    }
  }

  depends_on = [
    azapi_resource_action.add_accesspolicy_secret
  ]
}`, r.keyVaultSecretTemplate(data), data.RandomString)
}

func (r DataPlaneResource) keyVaultSecretTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azapi_client_config" "current" {}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        family = "A"
        name   = "standard"
      }
      tenantId                  = data.azapi_client_config.current.tenant_id
      enableSoftDelete          = true
      softDeleteRetentionInDays = 7
      enablePurgeProtection     = true
      accessPolicies            = []
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [body.properties.accessPolicies]
  }
  response_export_values = {
    "vaultUri" = "properties.vaultUri" 
  }
}

resource "azapi_resource_action" "add_accesspolicy_secret" {
  type        = "Microsoft.KeyVault/vaults/accessPolicies@2023-02-01"
  resource_id = "${azapi_resource.vault.id}/accessPolicies/add"
  method      = "PUT"
  body = {
    properties = {
      accessPolicies = [{
        tenantId = data.azapi_client_config.current.tenant_id
        objectId = data.azapi_client_config.current.object_id
        permissions = {
          secrets = [
            "Get", "List", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"
          ]
        }
      }]
    }
  }
}

`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) iotAppsUser(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "iotApp" {
  type      = "Microsoft.IoTCentral/iotApps@2021-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    sku = {
      name = "ST2"
    }
    properties = {
	  displayName = "acctest%[2]s"
      subdomain = "acctest%[2]s"
    }
  }
  response_export_values = {
    subdomain = "properties.subdomain"
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.IoTCentral/IoTApps/users@2022-07-31"
  parent_id = "${azapi_resource.iotApp.output.subdomain}.azureiotcentral.com"
  name      = "acctest%[2]s"
  body = {
    type = "email"
    roles = [
      {
        role = "ae2c9854-393b-4f97-8c42-479d70ce626e"
      }
    ]
    email = "user5@contoso.com"
  }
}

`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) timeouts(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "appconf" {
  type      = "Microsoft.AppConfiguration/configurationStores@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    sku = {
      name = "standard"
    }
  }
  response_export_values = {
    endpoint = "properties.endpoint"
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    appConfigDataOwnerRoleId = "value[?properties.roleName == 'App Configuration Data Owner'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.appconf.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.appConfigDataOwnerRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.AppConfiguration/configurationStores/keyValues@1.0"
  parent_id = replace(azapi_resource.appconf.output.endpoint, "https://", "")
  name      = "mykey"
  body = {
    content_type = ""
    value        = "myvalue"
  }

  depends_on = [
    azapi_resource.roleAssignment,
  ]

  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
    read   = "10m"
  }
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) oldConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = jsonencode({
    properties = {
      publicNetworkAccess = "Enabled"
    }
  })
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Purview/accounts/Account/collections@2019-11-01-preview"
  parent_id = "${azapi_resource.account.name}.purview.azure.com"
  name      = "defaultResourceSetRuleConfig"
  body = jsonencode({
    friendlyName = "Finance"
  })
  response_export_values = ["*"]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) headers(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "appconf" {
  type      = "Microsoft.AppConfiguration/configurationStores@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    sku = {
      name = "standard"
    }
  }
  response_export_values = {
    endpoint = "properties.endpoint"
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    appConfigDataOwnerRoleId = "value[?properties.roleName == 'App Configuration Data Owner'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.appconf.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.appConfigDataOwnerRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.AppConfiguration/configurationStores/keyValues@1.0"
  parent_id = replace(azapi_resource.appconf.output.endpoint, "https://", "")
  name      = "mykey"
  body = {
    content_type = ""
    value        = "myvalue"
  }

  create_headers = {
    "header1" = "create-value"
  }
  update_headers = {
    "header2" = "update-value"
  }
  delete_headers = {
    "header3" = "delete-value"
  }
  read_headers = {
    "header4" = "read-value"
  }

  depends_on = [
    azapi_resource.roleAssignment,
  ]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) queryParameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "appconf" {
  type      = "Microsoft.AppConfiguration/configurationStores@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    sku = {
      name = "standard"
    }
  }
  response_export_values = {
    endpoint = "properties.endpoint"
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    appConfigDataOwnerRoleId = "value[?properties.roleName == 'App Configuration Data Owner'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.appconf.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.appConfigDataOwnerRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.AppConfiguration/configurationStores/keyValues@1.0"
  parent_id = replace(azapi_resource.appconf.output.endpoint, "https://", "")
  name      = "mykey"
  body = {
    content_type = ""
    value        = "myvalue"
  }

  create_query_parameters = {
    "query1" = ["create-value"]
  }
  update_query_parameters = {
    "query1" = ["update-value"]
  }
  delete_query_parameters = {
    "query1" = ["delete-value"]
  }
  read_query_parameters = {
    "query1" = ["read-value"]
  }

  depends_on = [
    azapi_resource.roleAssignment,
  ]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) replaceTriggeredByExternalValues(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "appconf" {
  type      = "Microsoft.AppConfiguration/configurationStores@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    sku = {
      name = "standard"
    }
  }
  response_export_values = {
    endpoint = "properties.endpoint"
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    appConfigDataOwnerRoleId = "value[?properties.roleName == 'App Configuration Data Owner'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.appconf.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.appConfigDataOwnerRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.AppConfiguration/configurationStores/keyValues@1.0"
  parent_id = replace(azapi_resource.appconf.output.endpoint, "https://", "")
  name      = "mykey"
  body = {
    content_type = ""
    value        = "myvalue"
  }

  replace_triggers_external_values = [
    "value1"
  ]

  depends_on = [
    azapi_resource.roleAssignment,
  ]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) searchServiceIndex(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "searchService" {
  type      = "Microsoft.Search/searchServices@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      replicaCount   = 1
      partitionCount = 1
      hostingMode    = "default"
      authOptions = {
        aadOrApiKey = {
          aadAuthFailureMode = "http401WithBearerChallenge"
        }
      }
    }
    sku = {
      name = "basic"
    }
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    searchIndexDataContributorRoleId = "value[?properties.roleName == 'Search Index Data Contributor'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.searchService.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.searchIndexDataContributorRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Search/searchServices/indexes@2024-07-01"
  parent_id = "${azapi_resource.searchService.name}.search.windows.net"
  name      = "hotels-index"
  body = {
    fields = [
      {
        name       = "hotelId"
        type       = "Edm.String"
        key        = true
        searchable = false
      },
      {
        name       = "hotelName"
        type       = "Edm.String"
        searchable = true
      },
      {
        name       = "description"
        type       = "Edm.String"
        searchable = true
        analyzer   = "en.lucene"
      },
      {
        name       = "category"
        type       = "Edm.String"
        searchable = true
        filterable = true
      }
    ]
  }

  depends_on = [
    azapi_resource.roleAssignment,
  ]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) searchServiceDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "searchService" {
  type      = "Microsoft.Search/searchServices@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      replicaCount   = 1
      partitionCount = 1
      hostingMode    = "default"
      authOptions = {
        aadOrApiKey = {
          aadAuthFailureMode = "http401WithBearerChallenge"
        }
      }
    }
    sku = {
      name = "basic"
    }
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    searchIndexDataContributorRoleId = "value[?properties.roleName == 'Search Index Data Contributor'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.searchService.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.searchIndexDataContributorRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier = "Hot"
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  response_export_values = ["*"]
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2023-01-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "storageContainer" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-01-01"
  parent_id = "${azapi_resource.storageAccount.id}/blobServices/default"
  name      = "content"
  body = {
    properties = {}
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Search/searchServices/datasources@2024-07-01"
  parent_id = "${azapi_resource.searchService.name}.search.windows.net"
  name      = "mydatasource"
  body = {
    type        = "azureblob"
    credentials = {
      connectionString = "DefaultEndpointsProtocol=https;AccountName=${azapi_resource.storageAccount.name};AccountKey=${data.azapi_resource_action.listKeys.output.keys[0].value};EndpointSuffix=core.windows.net"
    }
    container = {
      name = azapi_resource.storageContainer.name
    }
  }

  depends_on = [
    azapi_resource.roleAssignment,
  ]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) searchServiceIndexer(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "searchService" {
  type      = "Microsoft.Search/searchServices@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      replicaCount   = 1
      partitionCount = 1
      hostingMode    = "default"
      authOptions = {
        aadOrApiKey = {
          aadAuthFailureMode = "http401WithBearerChallenge"
        }
      }
    }
    sku = {
      name = "basic"
    }
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    searchIndexDataContributorRoleId = "value[?properties.roleName == 'Search Index Data Contributor'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.searchService.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.searchIndexDataContributorRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier = "Hot"
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  response_export_values = ["*"]
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2023-01-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "storageContainer" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-01-01"
  parent_id = "${azapi_resource.storageAccount.id}/blobServices/default"
  name      = "content"
  body = {
    properties = {}
  }
}

resource "azapi_data_plane_resource" "index" {
  type      = "Microsoft.Search/searchServices/indexes@2024-07-01"
  parent_id = "${azapi_resource.searchService.name}.search.windows.net"
  name      = "hotels-index"
  body = {
    fields = [
      {
        name       = "hotelId"
        type       = "Edm.String"
        key        = true
        searchable = false
      },
      {
        name       = "hotelName"
        type       = "Edm.String"
        searchable = true
      }
    ]
  }
}

resource "azapi_data_plane_resource" "datasource" {
  type      = "Microsoft.Search/searchServices/datasources@2024-07-01"
  parent_id = "${azapi_resource.searchService.name}.search.windows.net"
  name      = "mydatasource"
  body = {
    type        = "azureblob"
    credentials = {
      connectionString = "DefaultEndpointsProtocol=https;AccountName=${azapi_resource.storageAccount.name};AccountKey=${data.azapi_resource_action.listKeys.output.keys[0].value};EndpointSuffix=core.windows.net"
    }
    container = {
      name = azapi_resource.storageContainer.name
    }
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Search/searchServices/indexers@2024-07-01"
  parent_id = "${azapi_resource.searchService.name}.search.windows.net"
  name      = "myindexer"
  body = {
    dataSourceName  = azapi_data_plane_resource.datasource.name
    targetIndexName = azapi_data_plane_resource.index.name
    schedule = {
      interval = "PT2H"
    }
  }
  depends_on = [
    azapi_resource.roleAssignment,
    azapi_data_plane_resource.datasource,
    azapi_data_plane_resource.index
  ]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) searchServiceSynonymMap(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azapi_resource" "searchService" {
  type      = "Microsoft.Search/searchServices@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      replicaCount   = 1
      partitionCount = 1
      hostingMode    = "default"
      authOptions = {
        aadOrApiKey = {
          aadAuthFailureMode = "http401WithBearerChallenge"
        }
      }
    }
    sku = {
      name = "basic"
    }
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    searchIndexDataContributorRoleId = "value[?properties.roleName == 'Search Index Data Contributor'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.searchService.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.searchIndexDataContributorRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Search/searchServices/synonymmaps@2024-07-01"
  parent_id = "${azapi_resource.searchService.name}.search.windows.net"
  name      = "mysynonymmap"
  body = {
    format = "solr"
    synonyms = "hotel, motel\nairport, aerodrome"
  }

  depends_on = [
    azapi_resource.roleAssignment,
  ]
}
`, data.LocationPrimary, data.RandomString)
}
