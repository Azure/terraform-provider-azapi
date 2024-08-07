package services_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type GenericResource struct{}

func defaultIgnores() []string {
	return []string{"ignore_casing", "ignore_missing_property", "schema_validation_enabled", "body", "locks", "removing_special_chars", "output"}
}

var testCertRaw, _ = os.ReadFile(filepath.Join("testdata", "automation_certificate_test.pfx"))

var testCertBase64 = base64.StdEncoding.EncodeToString(testCertRaw)

func TestAccGenericResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_dynamicSchema(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dynamicSchema(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_invalidVersionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config:      r.basicInvalidVersion(data),
			ExpectError: regexp.MustCompile("400 Bad Request"),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config:   r.basic(data),
			PlanOnly: true,
		},
	})
}

func TestAccGenericResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccGenericResource_importWithApiVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importWithApiVersion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       r.ImportIdFunc,
			ImportStateVerifyIgnore: defaultIgnores(),
		},
	})
}

func TestAccGenericResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_completeBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	// the identity block now is optional only, because framework doesn't allow computed+optional block, so the identity block couldn't be synced when it's set in the `body`
	importIgnores := []string{"identity.#", "identity.0.%", "identity.0.identity_ids.#", "identity.0.identity_ids.0", "identity.0.principal_id", "identity.0.tenant_id", "identity.0.type"}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(append(importIgnores, defaultIgnores()...)...),
	})
}

func TestAccGenericResource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.identityNone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.identityUserAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.identitySystemAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_defaultTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultTag(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("default"),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.defaultTagOverrideInBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("override"),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.defaultTag(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("default"),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.defaultTagOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("override"),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_defaultsNotApplicable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultsNotApplicable(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags").DoesNotExist(),
				check.That(data.ResourceName).Key("location").DoesNotExist(),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_defaultLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultLocation(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.defaultLocationOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationSecondary)),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_defaultParentId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultParentId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parent_id").HasValue(fmt.Sprintf("/subscriptions/%s", subscriptionId)),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAccGenericResource_defaultsNaming(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultNaming(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("acctestdefaultNaming"),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.defaultNamingOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("hclNaming"),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_defaultsNamingPrefixAndSuffix(t *testing.T) {
	t.Skip(`The default naming prefix and suffix are not compatible with the framework, because the computed name is not the same as the config name, the framework will fail as it's an invalid plan'`)
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultNamingWithPrefixAndSuffix(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_subscriptionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subscriptionScope(data, os.Getenv("ARM_SUBSCRIPTION_ID")),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAccGenericResource_extensionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionScope(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAccGenericResource_ignoreMissingProperty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreMissingProperty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAccGenericResource_ignoreCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreCasing(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAccGenericResource_deleteLROEndsWithNotFoundError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.deleteLROEndsWithNotFoundError(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_locks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.locks(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAccGenericResource_secretsInAsterisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.secretsInAsterisks(data, clientId, clientSecret),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_ignoreChanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreChanges(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_ignoreChangesArray(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreChangesArray(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_nonstandardLRO(t *testing.T) {
	t.Skip("This test is passing locally but failing in CI. Skipping for now.")
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.nonstandardLRO(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_nullLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.nullLocation(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_unknownName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.unknownName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.timeouts(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (GenericResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resourceType := state.Attributes["type"]
	id, err := parse.ResourceIDWithResourceType(state.ID, resourceType)
	if err != nil {
		return nil, err
	}

	_, err = client.ResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion)
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

func (GenericResource) ImportIdFunc(tfState *terraform.State) (string, error) {
	state := tfState.RootModule().Resources["azapi_resource.test"].Primary
	resourceType := state.Attributes["type"]
	id, err := parse.ResourceIDWithResourceType(state.ID, resourceType)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?api-version=%s", id.AzureResourceId, id.ApiVersion), nil
}

func (r GenericResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_account" "test" {
  name                = "acctest%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts/certificates@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_automation_account.test.id

  body = jsonencode({
    properties = {
      base64Value = "%[3]s"
    }
  })
}
`, r.template(data), data.RandomString, testCertBase64)
}

func (r GenericResource) dynamicSchema(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_account" "test" {
  name                = "acctest%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts/certificates@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_automation_account.test.id

  body = {
    properties = {
      base64Value = "%[3]s"
    }
  }
}
`, r.template(data), data.RandomString, testCertBase64)
}

func (r GenericResource) basicInvalidVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_account" "test" {
  name                = "acctest%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_resource" "test" {
  type                      = "Microsoft.Automation/automationAccounts/certificates@1999-01-01"
  name                      = "acctest%[2]s"
  parent_id                 = azurerm_automation_account.test.id
  schema_validation_enabled = false
  body = jsonencode({
    properties = {
      base64Value = "%[3]s"
    }
  })
}
`, r.template(data), data.RandomString, testCertBase64)
}

func (r GenericResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "import" {
  type      = azapi_resource.test.type
  name      = azapi_resource.test.name
  parent_id = azapi_resource.test.parent_id
  body = jsonencode({
    properties = {
      base64Value = "%s"
    }
  })
}
`, r.basic(data), testCertBase64)
}

func (r GenericResource) importWithApiVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_account" "test" {
  name                = "acctest%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts/certificates@2020-01-13-preview"
  name      = "acctest%[2]s"
  parent_id = azurerm_automation_account.test.id

  body = jsonencode({
    properties = {
      base64Value = "%[3]s"
    }
  })
}
`, r.template(data), data.RandomString, testCertBase64)
}

func (r GenericResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = "%[3]s"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })

  tags = {
    "Key" = "Value"
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) completeBody(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azapi" {
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azapi_resource" "test" {
  name                      = "acctest%[2]s"
  parent_id                 = azurerm_resource_group.test.id
  type                      = "Microsoft.Automation/automationAccounts@2023-11-01"
  schema_validation_enabled = false
  body                      = <<BODY
    {
      "location": "${azurerm_resource_group.test.location}",
      "identity": {
		"type": "SystemAssigned, UserAssigned",
        "userAssignedIdentities": {
          "${azurerm_user_assigned_identity.test.id}": {}
        }
      },
      "properties": {
        "sku": {
          "name": "Basic"
        }
      },
      "tags": {
        "key":"value"
      }
    }
  BODY
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) identityNone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = "%[3]s"

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = "%[3]s"
  identity {
    type = "SystemAssigned"
  }
  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = "%[3]s"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })

  tags = {
    "Key" = "Value"
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) defaultTag(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azapi" {
  default_tags = {
    key = "default"
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })

}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) defaultTagOverrideInBody(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azapi" {
  default_tags = {
    key = "default"
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
    tags = {
      key = "override"
    }
  })

}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) defaultTagOverrideInHcl(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azapi" {
  default_tags = {
    key = "default"
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })

  tags = {
    key = "override"
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) defaultLocation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azapi" {
  default_location = "%[3]s"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) defaultLocationOverrideInHcl(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azapi" {
  default_location = "%[3]s"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = "%[4]s"
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
    tags = {
      key = "override"
    }
  })

}
`, r.template(data), data.RandomString, data.LocationPrimary, data.LocationSecondary)
}

func (r GenericResource) defaultParentId(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azapi" {
}

resource "azapi_resource" "test" {
  type     = "Microsoft.Resources/resourceGroups@2023-07-01"
  name     = "acctest-%[2]d"
  location = "%[1]s"
}
`, data.LocationPrimary, data.RandomInteger)
}

func (r GenericResource) defaultNaming(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azapi" {
  default_name = "acctestdefaultNaming"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  parent_id = azurerm_resource_group.test.id

  location = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
}
`, r.template(data))
}

func (r GenericResource) defaultNamingOverrideInHcl(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azapi" {
  default_name = "acctestdefaultNaming"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "hclNaming"
  parent_id = azurerm_resource_group.test.id

  location = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
}
`, r.template(data))
}

func (r GenericResource) defaultNamingWithPrefixAndSuffix(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azapi" {
  default_naming_prefix = "p%[2]s-"
  default_naming_suffix = "-s%[2]s"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acc"
  parent_id = azurerm_resource_group.test.id

  location = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) defaultsNotApplicable(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azapi" {
  default_tags = {
    key = "default"
  }
  default_location = "%[3]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts/certificates@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_automation_account.test.id

  body = jsonencode({
    properties = {
      base64Value = "%[4]s"
    }
  })
}
`, r.template(data), data.RandomString, data.LocationPrimary, testCertBase64)
}

func (GenericResource) subscriptionScope(data acceptance.TestData, subscriptionId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Resources/resourceGroups@2023-07-01"
  name      = "acctestRG-%[1]d"
  parent_id = "/subscriptions/%[2]s"

  location = "%[3]s"
}
`, data.RandomInteger, subscriptionId, data.LocationPrimary)
}

// nolint staticcheck
func (r GenericResource) extensionScope(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azurerm_resource_group.test.id
  name      = "acctest-oi-%[2]d"
  location  = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  })
}

resource "azapi_resource" "onboardingState" {
  type      = "Microsoft.SecurityInsights/onboardingStates@2022-11-01"
  parent_id = azapi_resource.workspace.id
  name      = "default"
  body = jsonencode({
    properties = {
      customerManagedKey = false
    }
  })
}

resource "azapi_resource" "test" {
  type      = "Microsoft.SecurityInsights/watchlists@2022-11-01"
  parent_id = azapi_resource.workspace.id
  name      = "acctest-wl-%[2]d"
  body = jsonencode({
    properties = {
      displayName    = "test"
      itemsSearchKey = "k1"
      provider       = "Microsoft"
      source         = ""
    }
  })
  depends_on = [azapi_resource.onboardingState]
}
`, r.template(data), data.RandomInteger)
}

func (r GenericResource) ignoreMissingProperty(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azapi_resource" "test" {
  type      = "Microsoft.AppPlatform/spring/storages@2023-12-01"
  name      = "acctest-ss-%[2]d"
  parent_id = azurerm_spring_cloud_service.test.id

  body = jsonencode({
    properties = {
      accountKey  = azurerm_storage_account.test.primary_access_key
      accountName = azurerm_storage_account.test.name
      storageType = "StorageAccount"
    }
  })

  ignore_missing_property = true
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r GenericResource) ignoreCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azapi_resource" "test" {
  type      = "Microsoft.AppPlatform/spring/storages@2024-05-01-preview"
  name      = "acctest-ss-%[2]d"
  parent_id = azurerm_spring_cloud_service.test.id

  body = jsonencode({
    properties = {
      accountKey  = azurerm_storage_account.test.primary_access_key
      accountName = azurerm_storage_account.test.name
      storageType = "storageaccount"
    }
  })

  schema_validation_enabled = false
  ignore_casing             = true
  ignore_missing_property   = true
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r GenericResource) deleteLROEndsWithNotFoundError(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "test" {
  type      = "Microsoft.ServiceBus/namespaces@2022-10-01-preview"
  name      = "acctest-sb-%[2]d"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
  body = jsonencode({
    sku = {
      name = "Premium"
    }
  })
}

`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r GenericResource) locks(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_route_table" "test" {
  name                = "acctestrt%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/routeTables/routes@2023-09-01"
  name      = "first%[2]d"
  parent_id = azurerm_route_table.test.id
  body = jsonencode({
    properties = {
      nextHopType   = "VnetLocal"
      addressPrefix = "10.1.0.0/16"
    }
  })

  locks = [azurerm_route_table.test.id, azurerm_resource_group.test.id]
}

resource "azapi_resource" "test2" {
  type      = "Microsoft.Network/routeTables/routes@2023-09-01"
  name      = "second%[2]d"
  parent_id = azurerm_route_table.test.id
  body = jsonencode({
    properties = {
      nextHopType   = "VnetLocal"
      addressPrefix = "10.3.0.0/16"
    }
  })

  locks = [azurerm_route_table.test.id, azurerm_resource_group.test.id]
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r GenericResource) secretsInAsterisks(data acceptance.TestData, clientId, clientSecret string) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_gateway" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
}

resource "azapi_resource" "test" {
  type      = "Microsoft.AppPlatform/Spring/apiPortals@2022-12-01"
  parent_id = azurerm_spring_cloud_service.test.id
  name      = "default"
  body = jsonencode({
    properties = {
      gatewayIds = [azurerm_spring_cloud_gateway.test.id]
      httpsOnly  = false
      public     = false
      ssoProperties = {
        clientId     = "%[4]s"
        clientSecret = "%[5]s"
        issuerUri    = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0"
        scope        = ["read"]
      }
    }
  })
  ignore_casing = true
}
`, r.template(data), data.RandomInteger, data.RandomString, clientId, clientSecret)
}

func (r GenericResource) ignoreChanges(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]d"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      sku = {
        name = "Free"
      }
    }
  })

  ignore_body_changes = ["properties.sku.name"]
}


`, r.template(data), data.RandomInteger)
}

func (r GenericResource) ignoreChangesArray(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azurerm_resource_group.test.id
  name      = "acctest%[2]d"
  location  = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/16",
        ]
      }
      subnets = [
        {
          name = "first"
          properties = {
            addressPrefix = "10.0.1.0/24"
          }
        }
      ]
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
  ignore_body_changes       = ["properties.subnets"]
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.test.id
  name      = "second"
  body = jsonencode({
    properties = {
      addressPrefix = "10.0.2.0/24"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
`, r.template(data), data.RandomInteger)
}

func (r GenericResource) nonstandardLRO(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                 = "acctestsc%[2]s"
  storage_account_name = azurerm_storage_account.test.name
}


resource "azapi_resource" "test" {
  type      = "Microsoft.CostManagement/exports@2022-10-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  body = jsonencode({
    properties = {
      schedule = {
        recurrence = "Monthly"
        recurrencePeriod = {
          from = "2030-12-29T00:00:00Z"
          to   = "2030-12-30T00:00:00Z"
        }
        status = "Inactive"
      }
      definition = {
        timeframe = "TheLastMonth"
        type      = "Usage"

      }
      format = "Csv"
      deliveryInfo = {
        destination = {
          rootFolderPath = "test"
          container      = azurerm_storage_container.test.name
          resourceId     = azurerm_storage_account.test.id
        }
      }
    }
  })
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) nullLocation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_application_insights" "test" {
  name                = "accappinsights%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acckeyvault%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[2]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctestmlws%[2]s"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
  public_network_access_enabled = true
  managed_network {
    isolation_mode = "AllowOnlyApprovedOutbound"
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.MachineLearningServices/workspaces/outboundRules@2023-10-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_machine_learning_workspace.test.id
  body = jsonencode({
    properties = {
      category    = "UserDefined"
      status      = "Active"
      type        = "FQDN"
      destination = "example.org"
    }
  })
  locks = [azurerm_machine_learning_workspace.test.id]
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) unknownName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {}

resource "random_string" "suffix" {
  length  = 3
  special = false
  upper   = false
}

resource "azapi_resource" "test" {
  type      = "Microsoft.KeyVault/vaults@2023-07-01"
  parent_id = azurerm_resource_group.test.id
  name      = "acctest${random_string.suffix.result}"
  location  = azurerm_resource_group.test.location
  body = {
    properties = {
      accessPolicies = [
      ]
      createMode                   = "default"
      enableRbacAuthorization      = false
      enableSoftDelete             = true
      enabledForDeployment         = false
      enabledForDiskEncryption     = false
      enabledForTemplateDeployment = false
      publicNetworkAccess          = "Enabled"
      sku = {
        family = "A"
        name   = "standard"
      }
      softDeleteRetentionInDays = 7
      tenantId                  = data.azurerm_client_config.current.tenant_id
    }
  }
}
`, r.template(data))
}

func (r GenericResource) timeouts(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
    read   = "10m"
  }
}
`, r.template(data), data.RandomString)
}

func (GenericResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString)
}
