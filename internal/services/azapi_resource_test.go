package services_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type GenericResource struct{}

func defaultIgnores() []string {
	return []string{"ignore_casing", "ignore_missing_property", "schema_validation_enabled", "body", "locks", "output", "create_", "delete_", "update_", "read_"}
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
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAccGenericResource_resourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "resourceGroup")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// Template contains only a Resource Group
			Config: r.template(data),
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
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config:      r.basicInvalidVersion(data),
			ExpectError: regexp.MustCompile("400 Bad Request"),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
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
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
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
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, append(importIgnores, defaultIgnores()...)...),
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
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config: r.identityUserAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config: r.identitySystemAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
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
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config: r.defaultTagOverrideInBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("override"),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config: r.defaultTag(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("default"),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config: r.defaultTagOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("override"),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
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
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
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
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config: r.defaultLocationOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationSecondary)),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
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
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config: r.defaultNamingOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("hclNaming"),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
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
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
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

func TestAccGenericResource_nonstandardLRO(t *testing.T) {
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
			Config:            r.nullLocation(data),
			ExternalProviders: externalProvidersAzurerm(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_computedLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:            r.computedLocation(data),
			ExternalProviders: externalProvidersAzurerm(),
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
			ExternalProviders: map[string]resource.ExternalProvider{
				"random": {
					Source:            "registry.terraform.io/hashicorp/random",
					VersionConstraint: "= 3.6.1",
				},
			},
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_unknownNameWithSensitiveBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.unknownNameWithSensitiveBody(data),
			ExternalProviders: map[string]resource.ExternalProvider{
				"random": {
					Source:            "registry.terraform.io/hashicorp/random",
					VersionConstraint: "= 3.6.1",
				},
			},
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

func TestAccGenericResource_replaceTriggeredBy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.replaceTriggeredByValue1(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.replaceTriggeredByValue2(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.replaceTriggeredByValueNull(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_withRetry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withRetry(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.headers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_queryParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.queryParameters(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_replaceTriggersRefs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.replaceTriggersRefs(data, "S0"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.replaceTriggersRefs(data, "E0"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_defaultOutput(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultOutput(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.properties.automationHybridServiceUrl").Exists(),
			),
		},
	})
}

func TestAccGenericResource_unknownDiscriminator(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.unknownDiscriminator(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
			ExternalProviders:  knownExternalProvidersAzurerm(),
		},
	})
}

func TestAccGenericResource_moveResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:            r.moveResourceSetup(data),
			Check:             resource.ComposeTestCheckFunc(),
			ExternalProviders: externalProvidersAzurerm(),
		},
		{
			Config: r.moveResourceStartMoving(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExternalProviders: externalProvidersAzurerm(),
		},
		{
			Config: r.moveResourceUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExternalProviders: externalProvidersAzurerm(),
		},
	})
}

func TestAccGenericResource_moveStorageContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:            r.moveStorageContainerSetup(data),
			ExternalProviders: externalProvidersAzurerm(),
		},
		{
			Config: r.moveStorageContainerStartMoving(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExternalProviders: externalProvidersAzurerm(),
		},
		{
			Config: r.moveStorageContainerUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExternalProviders: externalProvidersAzurerm(),
		},
	})
}

func TestAccGenericResource_moveKeyVaultSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:            r.moveKeyVaultSecretSetup(data),
			ExternalProviders: externalProvidersAzurerm(),
		},
		{
			Config: r.moveKeyVaultSecretStartMoving(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExternalProviders: externalProvidersAzurerm(),
		},
		{
			Config: r.moveKeyVaultSecretUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExternalProviders: externalProvidersAzurerm(),
		},
		{
			Config: r.moveKeyVaultSecretRemoved(data),
			Check:  resource.ComposeTestCheckFunc(
			// resource should be removed from state; no existence check
			),
			ExternalProviders: externalProvidersAzurerm(),
		},
	})
}

func TestAccGenericResource_moveKeyVaultKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:            r.moveKeyVaultKeySetup(data),
			ExternalProviders: externalProvidersAzurerm(),
		},
		{
			Config: r.moveKeyVaultKeyStartMoving(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExternalProviders: externalProvidersAzurerm(),
		},
		{
			Config: r.moveKeyVaultKeyUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExternalProviders: externalProvidersAzurerm(),
		},
		{
			Config:            r.moveKeyVaultKeyRemoved(data),
			Check:             resource.ComposeTestCheckFunc(),
			ExternalProviders: externalProvidersAzurerm(),
		},
	})
}

func TestAccGenericResource_SensitiveBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.SensitiveBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAccGenericResource_SensitiveBodyVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	ignores := defaultIgnores()
	ignores = append(ignores, "tags")
	ignores = append(ignores, "sensitive_body_version")
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.SensitiveBodyWithHash(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.tags.tag1").HasValue("tag1-value"),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, ignores...),
		{
			Config: r.SensitiveBodyWithVersion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.tags.tag1").HasValue("tag1-value"),
				check.That(data.ResourceName).Key("output.tags.tag2").HasValue("tag2-value2"),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, ignores...),
		{
			Config: r.SensitiveBodyWithVersionMultipleTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.tags.tag1").HasValue("tag1-value"),
				check.That(data.ResourceName).Key("output.tags.tag2").HasValue("tag2-value2"),
				check.That(data.ResourceName).Key("output.tags.tag3").DoesNotExist(),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, ignores...),
		{
			Config: r.SensitiveBodyWithHashMultipleTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.tags.tag1").HasValue("tag1-value"),
				check.That(data.ResourceName).Key("output.tags.tag2").HasValue("tag2-value3"),
				check.That(data.ResourceName).Key("output.tags.tag3").HasValue("tag3-value"),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, ignores...),
	})
}

func TestAccGenericResource_multipleIdentityIds(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleIdentityIds(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.multipleIdentityIds(data, true),
			ExpectNonEmptyPlan: false,
		},
		{
			Config:             r.multipleIdentityIds(data, true),
			ExpectNonEmptyPlan: false,
		},
	})
}

func TestAccGenericResource_BodySemanticallyEqualToRemote(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.automationAccountBasic(data, "2023-11-01"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.automationAccountBasic(data, "2024-10-23"),
			ExpectNonEmptyPlan: false,
		},
		{
			Config:             r.automationAccountComplete(data, "2024-10-23"),
			ExpectNonEmptyPlan: false,
		},
		{
			Config:             r.automationAccountComplete(data, "2023-11-01"),
			ExpectNonEmptyPlan: false,
		},
		{
			Config:             r.automationAccountCompleteStrictChangeDetection(data, "2023-11-01"),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_IgnoreNullProperty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.automationAccountCompleteWithNullPropertiesSetup(data, "2024-10-23"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.automationAccountCompleteWithNullProperties(data, "2024-10-23"),
			ExpectNonEmptyPlan: false,
		},
	})
}

func TestAccGenericResource_MovingFromAzureRM(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:            r.automationAccountAzureRM(data),
			ExternalProviders: knownExternalProvidersAzurerm(),
		},
		{
			Config:             r.automationAccountAzureRMMovedBasic(data, "2023-11-01"),
			ExpectNonEmptyPlan: false,
		},
		{
			Config:             r.automationAccountAzureRMMovedBasic(data, "2024-10-23"),
			ExternalProviders:  knownExternalProvidersAzurerm(),
			ExpectNonEmptyPlan: false,
		},
		{
			Config:             r.automationAccountAzureRMMovedComplete(data, "2024-10-23"),
			ExternalProviders:  knownExternalProvidersAzurerm(),
			ExpectNonEmptyPlan: false,
		},
		{
			Config:             r.automationAccountAzureRMMovedComplete(data, "2023-11-01"),
			ExternalProviders:  knownExternalProvidersAzurerm(),
			ExpectNonEmptyPlan: false,
		},
		{
			Config:             r.automationAccountAzureRMMovedCompleteStrictChangeDetection(data, "2023-11-01"),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_modifyPlanSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.modifyPlanSubnet(data),
		},
		{
			Config:             r.modifyPlanSubnetUpdate(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_modifyPlanAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.modifyPlanAccount(data),
		},
		{
			Config:             r.modifyPlanAccountUpdate(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func (GenericResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resourceType := state.Attributes["type"]
	id, err := parse.ResourceIDWithResourceType(state.ID, resourceType)
	if err != nil {
		return nil, err
	}

	_, err = client.ResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion, clients.DefaultRequestOptions())
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

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
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

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts/certificates@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.automationAccount.id

  body = {
    properties = {
      base64Value = "%[3]s"
    }
  }
}
`, r.template(data), data.RandomString, testCertBase64)
}

func (r GenericResource) withRetry(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
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

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts/certificates@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.automationAccount.id

  retry = {
    error_message_regex = ["test error"]
  }

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

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
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


resource "azapi_resource" "test" {
  type                      = "Microsoft.Automation/automationAccounts/certificates@1999-01-01"
  name                      = "acctest%[2]s"
  parent_id                 = azapi_resource.automationAccount.id
  schema_validation_enabled = false
  body = {
    properties = {
      base64Value = "%[3]s"
    }
  }
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
  body = {
    properties = {
      base64Value = "%s"
    }
  }
}
`, r.basic(data), testCertBase64)
}

func (r GenericResource) importWithApiVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
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

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts/certificates@2020-01-13-preview"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.automationAccount.id

  body = {
    properties = {
      base64Value = "%[3]s"
    }
  }
}
`, r.template(data), data.RandomString, testCertBase64)
}

func (r GenericResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id

  location = "%[3]s"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }

  tags = {
    "Key" = "Value"
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) completeJsonBody(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "test" {
  name                      = "acctest%[2]s"
  parent_id                 = azapi_resource.resourceGroup.id
  type                      = "Microsoft.Automation/automationAccounts@2023-11-01"
  schema_validation_enabled = false
  body = jsonencode({
    location = azapi_resource.resourceGroup.location
    identity = {
      type = "SystemAssigned, UserAssigned"
      userAssignedIdentities = {
        (azapi_resource.userAssignedIdentity.id) = {}
      }
    }
    properties = {
      sku = {
        name = "Basic"
      }
    }
    tags = {
      "Key" = "Value"
    }
  })
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) completeBody(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "test" {
  name                      = "acctest%[2]s"
  parent_id                 = azapi_resource.resourceGroup.id
  type                      = "Microsoft.Automation/automationAccounts@2023-11-01"
  schema_validation_enabled = false
  body = {
    location = azapi_resource.resourceGroup.location
    identity = {
      type = "SystemAssigned, UserAssigned"
      userAssignedIdentities = {
        (azapi_resource.userAssignedIdentity.id) = {}
      }
    }
    properties = {
      sku = {
        name = "Basic"
      }
    }
    tags = {
      "Key" = "Value"
    }
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) identityNone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id

  location = "%[3]s"

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id

  location = "%[3]s"
  identity {
    type = "SystemAssigned"
  }
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id

  location = "%[3]s"
  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }

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
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }

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
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
    tags = {
      key = "override"
    }
  }

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
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }

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
  parent_id = azapi_resource.resourceGroup.id

  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
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
  parent_id = azapi_resource.resourceGroup.id

  location = "%[4]s"
  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
    tags = {
      key = "override"
    }
  }

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
`, r.template(data))
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

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
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

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts/certificates@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.automationAccount.id

  body = {
    properties = {
      base64Value = "%[4]s"
    }
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary, testCertBase64)
}

func (GenericResource) subscriptionScope(data acceptance.TestData, subscriptionId string) string {
	return fmt.Sprintf(`

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
  name      = "acctest-oi-%[2]d"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
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
  }
}

resource "azapi_resource" "onboardingState" {
  type      = "Microsoft.SecurityInsights/onboardingStates@2022-11-01"
  parent_id = azapi_resource.workspace.id
  name      = "default"
  body = {
    properties = {
      customerManagedKey = false
    }
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.SecurityInsights/watchlists@2022-11-01"
  parent_id = azapi_resource.workspace.id
  name      = "acctest-wl-%[2]d"
  body = {
    properties = {
      displayName    = "test"
      itemsSearchKey = "k1"
      provider       = "Microsoft"
      source         = ""
    }
  }
  depends_on = [azapi_resource.onboardingState]
}
`, r.template(data), data.RandomInteger)
}

func (r GenericResource) ignoreMissingProperty(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestsa%[3]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = false
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      isHnsEnabled                 = false
      isNfsV3Enabled               = false
      isSftpEnabled                = false
      minimumTlsVersion            = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2023-05-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "Spring" {
  type      = "Microsoft.AppPlatform/Spring@2024-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-sc-%[2]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      zoneRedundant = false
    }
    sku = {
      name = "S0"
    }
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.AppPlatform/Spring/storages@2023-12-01"
  name      = "acctest-ss-%[2]d"
  parent_id = azapi_resource.Spring.id

  body = {
    properties = {
      accountKey  = try(data.azapi_resource_action.listKeys.output.keys[0].value, jsondecode(data.azapi_resource_action.listKeys.output).keys[0].value)
      accountName = azapi_resource.storageAccount.name
      storageType = "StorageAccount"
    }
  }

  ignore_missing_property = true
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r GenericResource) ignoreCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestsa%[3]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = false
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      isHnsEnabled                 = false
      isNfsV3Enabled               = false
      isSftpEnabled                = false
      minimumTlsVersion            = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

data "azapi_resource_action" "listKeys" {
  type                   = "Microsoft.Storage/storageAccounts@2023-05-01"
  resource_id            = azapi_resource.storageAccount.id
  action                 = "listKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "Spring" {
  type      = "Microsoft.AppPlatform/Spring@2024-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-sc-%[2]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      zoneRedundant = false
    }
    sku = {
      name = "S0"
    }
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.AppPlatform/Spring/storages@2024-05-01-preview"
  name      = "acctest-ss-%[2]d"
  parent_id = azapi_resource.Spring.id

  body = {
    properties = {
      accountKey  = try(data.azapi_resource_action.listKeys.output.keys[0].value, jsondecode(data.azapi_resource_action.listKeys.output).keys[0].value)
      accountName = azapi_resource.storageAccount.name
      storageType = "storageaccount"
    }
  }

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
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    sku = {
      name = "Premium"
    }
  }
}

`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r GenericResource) locks(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "routeTable" {
  type      = "Microsoft.Network/routeTables@2024-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestrt%[2]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      disableBgpRoutePropagation = false
    }
  }
  lifecycle {
    ignore_changes = [body.properties.routes]
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/routeTables/routes@2023-09-01"
  name      = "first%[2]d"
  parent_id = azapi_resource.routeTable.id
  body = {
    properties = {
      nextHopType   = "VnetLocal"
      addressPrefix = "10.1.0.0/16"
    }
  }

  locks = [azapi_resource.routeTable.id, azapi_resource.resourceGroup.id]
}

resource "azapi_resource" "test2" {
  type      = "Microsoft.Network/routeTables/routes@2023-09-01"
  name      = "second%[2]d"
  parent_id = azapi_resource.routeTable.id
  body = {
    properties = {
      nextHopType   = "VnetLocal"
      addressPrefix = "10.3.0.0/16"
    }
  }

  locks = [azapi_resource.resourceGroup.id, azapi_resource.routeTable.id]
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (GenericResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2023-07-01"
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString)
}

func (r GenericResource) secretsInAsterisks(data acceptance.TestData, clientId, clientSecret string) string {
	return fmt.Sprintf(`
%[1]s

data "azapi_client_config" "current" {
}

resource "azapi_resource" "Spring" {
  type      = "Microsoft.AppPlatform/Spring@2024-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-sc-%[2]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      zoneRedundant = false
    }
    sku = {
      name = "E0"
    }
  }
}

resource "azapi_resource" "gateway" {
  type      = "Microsoft.AppPlatform/Spring/gateways@2024-05-01-preview"
  parent_id = azapi_resource.Spring.id
  name      = "default"
  body = {
    properties = {
      httpsOnly = false
      public    = false
    }
    sku = {
      capacity = 1
      name     = "E0"
      tier     = "Enterprise"
    }
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.AppPlatform/Spring/apiPortals@2022-12-01"
  parent_id = azapi_resource.Spring.id
  name      = "default"
  body = {
    properties = {
      gatewayIds = [azapi_resource.gateway.id]
      httpsOnly  = false
      public     = false
      ssoProperties = {
        clientId     = "%[4]s"
        clientSecret = "%[5]s"
        issuerUri    = "https://login.microsoftonline.com/${data.azapi_client_config.current.tenant_id}/v2.0"
        scope        = ["read"]
      }
    }
  }
  ignore_casing = true
}
`, r.template(data), data.RandomInteger, data.RandomString, clientId, clientSecret)
}

func (r GenericResource) nonstandardLRO(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestsa%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = false
      allowCrossTenantReplication  = true
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      isHnsEnabled                 = false
      isNfsV3Enabled               = false
      isSftpEnabled                = false
      minimumTlsVersion            = "TLS1_2"
      networkAcls = {
        defaultAction = "Allow"
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

data "azapi_resource_id" "blobService" {
  type      = "Microsoft.Storage/storageAccounts/blobServices@2023-05-01"
  parent_id = azapi_resource.storageAccount.id
  name      = "default"
}

resource "azapi_resource" "container" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01"
  parent_id = data.azapi_resource_id.blobService.id
  name      = "acctestsc%[2]s"
  body = {
    properties = {
    }
  }
}


resource "azapi_resource" "test" {
  type      = "Microsoft.CostManagement/exports@2022-10-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  body = {
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
          container      = azapi_resource.container.name
          resourceId     = azapi_resource.storageAccount.id
        }
      }
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) nullLocation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

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

data "azurerm_client_config" "current" {}

resource "azurerm_application_insights" "test" {
  name                = "accappinsights%[2]s"
  location            = azapi_resource.resourceGroup.location
  resource_group_name = azapi_resource.resourceGroup.name
  application_type    = "web"

  lifecycle {
    ignore_changes = [workspace_id]
  }
}

resource "azurerm_key_vault" "test" {
  name                     = "acckeyvault%[2]s"
  location                 = azapi_resource.resourceGroup.location
  resource_group_name      = azapi_resource.resourceGroup.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestsa%[2]s"
  location                        = azapi_resource.resourceGroup.location
  resource_group_name             = azapi_resource.resourceGroup.name
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctestmlws%[2]s"
  location                = azapi_resource.resourceGroup.location
  resource_group_name     = azapi_resource.resourceGroup.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  managed_network {
    isolation_mode = "AllowOnlyApprovedOutbound"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.MachineLearningServices/workspaces/outboundRules@2023-10-01"
  name      = "acctest%[2]s"
  parent_id = azurerm_machine_learning_workspace.test.id
  body = {
    properties = {
      category    = "UserDefined"
      status      = "Active"
      type        = "FQDN"
      destination = "example.org"
    }
  }
  locks = [azurerm_machine_learning_workspace.test.id]
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) computedLocation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "namespace" {
  type      = "Microsoft.EventHub/namespaces@2022-01-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      disableLocalAuth     = false
      isAutoInflateEnabled = false
      publicNetworkAccess  = "Enabled"
    }
    sku = {
      capacity = 1
      name     = "Basic"
      tier     = "Basic"
    }
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.EventHub/namespaces/authorizationRules@2021-11-01"
  parent_id = azapi_resource.namespace.id
  name      = "acctest%[2]s"
  body = {
    properties = {
      rights = [
        "Listen",
        "Send",
        "Manage",
      ]
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) unknownName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azapi_client_config" "current" {}

resource "random_string" "suffix" {
  length  = 3
  special = false
  upper   = false
}

resource "azapi_resource" "test" {
  type      = "Microsoft.KeyVault/vaults@2023-07-01"
  name      = "acctest${random_string.suffix.result}"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
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
      tenantId                  = data.azapi_client_config.current.tenant_id
    }
  }
}
`, r.template(data))
}

func (r GenericResource) unknownNameWithSensitiveBody(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azapi_client_config" "current" {}

resource "random_string" "suffix" {
  length  = 3
  special = false
  upper   = false
}

resource "azapi_resource" "test" {
  type      = "Microsoft.KeyVault/vaults@2023-07-01"
  name      = "acctest${random_string.suffix.result}"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
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
      sku = {
        family = "A"
        name   = "standard"
      }
      softDeleteRetentionInDays = 7
      tenantId                  = data.azapi_client_config.current.tenant_id
    }
  }
  sensitive_body = {
    properties = {
      publicNetworkAccess = "Enabled"
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
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
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

func (r GenericResource) replaceTriggeredByValue1(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
  replace_triggers_external_values = [
    "value1"
  ]
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) replaceTriggeredByValue2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
  replace_triggers_external_values = [
    "value2"
  ]
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) replaceTriggeredByValueNull(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
  replace_triggers_external_values = null
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) headers(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "test" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
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
}`, data.RandomInteger, data.LocationPrimary)
}

func (r GenericResource) queryParameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "test" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
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
}`, data.RandomInteger, data.LocationPrimary)
}

func (r GenericResource) oldConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
  tags = {
    env = "prod"
  }
  response_export_values = ["properties"]
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) replaceTriggersRefs(data acceptance.TestData, skuName string) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.AppPlatform/Spring@2024-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-sc-%[2]d"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      zoneRedundant = false
    }
    sku = {
      name = "%[3]s"
    }
  }
  replace_triggers_refs = ["sku.name"]
}
`, r.template(data), data.RandomInteger, skuName)
}

func (r GenericResource) defaultOutput(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azapi" {
  disable_default_output = false
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
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
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveResourceSetup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_cognitive_account" "test" {
  name                = "acctest%[2]s"
  location            = azapi_resource.resourceGroup.location
  resource_group_name = azapi_resource.resourceGroup.name
  kind                = "Face"
  sku_name            = "S0"
  tags = {
    Acceptance = "Test"
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveResourceStartMoving(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

moved {
  from = azurerm_cognitive_account.test
  to   = azapi_resource.test
}

resource "azapi_resource" "test" {
  type      = "Microsoft.CognitiveServices/accounts@2024-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "Face"
    properties = {
      allowedFqdnList               = []
      disableLocalAuth              = false
      dynamicThrottlingEnabled      = false
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = false
    }
    sku = {
      name = "S0"
    }
  }
  tags = {
    Acceptance = "Test"
  }
  ignore_casing             = false
  schema_validation_enabled = true
  ignore_missing_property   = true
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveResourceUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

moved {
  from = azurerm_cognitive_account.test
  to   = azapi_resource.test
}

resource "azapi_resource" "test" {
  type      = "Microsoft.CognitiveServices/accounts@2024-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "Face"
    properties = {
      allowedFqdnList               = []
      disableLocalAuth              = false
      dynamicThrottlingEnabled      = false
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = false
      restore                       = null
    }
    sku = {
      name = "S0"
    }
  }
  tags = {
    Acceptance = "Test"
  }
  ignore_casing             = false
  schema_validation_enabled = true
  ignore_missing_property   = true
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveStorageContainerSetup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_storage_account" "sa" {
  name                     = "acctestsa%[2]s"
  location                 = azapi_resource.resourceGroup.location
  resource_group_name      = azapi_resource.resourceGroup.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "ct" {
  name                  = "acctestct%[2]s"
  storage_account_name  = azurerm_storage_account.sa.name
  container_access_type = "private"
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveStorageContainerStartMoving(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

moved {
  from = azurerm_storage_container.ct
  to   = azapi_resource.test
}

resource "azurerm_storage_account" "sa" {
  name                     = "acctestsa%[2]s"
  location                 = azapi_resource.resourceGroup.location
  resource_group_name      = azapi_resource.resourceGroup.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

//resource "azurerm_storage_container" "ct" {
//  name                  = "acctestct%[2]s"
//  storage_account_name  = azurerm_storage_account.sa.name
//  container_access_type = "private"
//}

data "azapi_resource_id" "blobService" {
  type      = "Microsoft.Storage/storageAccounts/blobServices@2023-05-01"
  parent_id = azurerm_storage_account.sa.id
  name      = "default"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01"
  parent_id = data.azapi_resource_id.blobService.id
  name      = "acctestct%[2]s"
  body = {
    properties = {}
  }
  ignore_casing             = false
  schema_validation_enabled = true
  ignore_missing_property   = true
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveStorageContainerUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

moved {
  from = azurerm_storage_container.ct
  to   = azapi_resource.test
}

resource "azurerm_storage_account" "sa" {
  name                     = "acctestsa%[2]s"
  location                 = azapi_resource.resourceGroup.location
  resource_group_name      = azapi_resource.resourceGroup.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

//resource "azurerm_storage_container" "ct" {
//  name                  = "acctestct%[2]s"
//  storage_account_name  = azurerm_storage_account.sa.name
//  container_access_type = "private"
//}

data "azapi_resource_id" "blobService" {
  type      = "Microsoft.Storage/storageAccounts/blobServices@2023-05-01"
  parent_id = azurerm_storage_account.sa.id
  name      = "default"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01"
  parent_id = data.azapi_resource_id.blobService.id
  name      = "acctestct%[2]s"
  body = {
    properties = {
      publicAccess = "Blob"
    }
  }
  ignore_casing             = false
  schema_validation_enabled = true
  ignore_missing_property   = true
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveKeyVaultSecretSetup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "kv" {
  name                       = "acctestkv%[2]s"
  location                   = azapi_resource.resourceGroup.location
  resource_group_name        = azapi_resource.resourceGroup.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
    ]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

resource "azurerm_key_vault_secret" "sec" {
  name         = "acctestsecret%[2]s"
  value        = "s3cr3tValue"
  key_vault_id = azurerm_key_vault.kv.id
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveKeyVaultSecretStartMoving(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

moved {
  from = azurerm_key_vault_secret.sec
  to   = azapi_resource.test
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "kv" {
  name                       = "acctestkv%[2]s"
  location                   = azapi_resource.resourceGroup.location
  resource_group_name        = azapi_resource.resourceGroup.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
    ]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

//resource "azurerm_key_vault_secret" "sec" {
//  name         = "acctestsecret%[2]s"
//  value        = "s3cr3tValue"
//  key_vault_id = azurerm_key_vault.kv.id
//}

resource "azapi_resource" "test" {
  type      = "Microsoft.KeyVault/vaults/secrets@2024-11-01"
  parent_id = azurerm_key_vault.kv.id
  name      = "acctestsecret%[2]s"
  body = {
    properties = {
      value = "s3cr3tValue"
    }
  }
  ignore_missing_property = true
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveKeyVaultSecretUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azurerm" {
  features {}
}

moved {
  from = azurerm_key_vault_secret.sec
  to   = azapi_resource.test
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "kv" {
  name                       = "acctestkv%[2]s"
  location                   = azapi_resource.resourceGroup.location
  resource_group_name        = azapi_resource.resourceGroup.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
    ]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

//resource "azurerm_key_vault_secret" "sec" {
//  name         = "acctestsecret%[2]s"
//  value        = "s3cr3tValue"
//  key_vault_id = azurerm_key_vault.kv.id
//}

resource "azapi_resource" "test" {
  type      = "Microsoft.KeyVault/vaults/secrets@2024-11-01"
  parent_id = azurerm_key_vault.kv.id
  name      = "acctestsecret%[2]s"
  body = {
    properties = {
      value = "updatedS3cr3tValue"
    }
  }
  ignore_missing_property = true
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveKeyVaultSecretRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azurerm" {
  features {}
}

moved {
  from = azurerm_key_vault_secret.sec
  to   = azapi_resource.test
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "kv" {
  name                       = "acctestkv%[2]s"
  location                   = azapi_resource.resourceGroup.location
  resource_group_name        = azapi_resource.resourceGroup.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
    ]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

//resource "azurerm_key_vault_secret" "sec" {
//  name         = "acctestsecret%[2]s"
//  value        = "s3cr3tValue"
//  key_vault_id = azurerm_key_vault.kv.id
//}

//resource "azapi_resource" "test" {
//  type      = "Microsoft.KeyVault/vaults/secrets@2024-11-01"
//  parent_id = azurerm_key_vault.kv.id
//  name      = "acctestsecret%[2]s"
//  body = {
//    properties = {
//      value = "updatedS3cr3tValue"
//    }
//  }
//  ignore_missing_property = true
//}

removed {
  from = azapi_resource.test
  lifecycle {
    destroy = false
  }
}

`, r.template(data), data.RandomString)
}

func (r GenericResource) moveKeyVaultKeySetup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "kv" {
  name                       = "acctestkv%[2]s"
  location                   = azapi_resource.resourceGroup.location
  resource_group_name        = azapi_resource.resourceGroup.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
      "Update",
      "GetRotationPolicy",
      "SetRotationPolicy",
    ]
  }
}

resource "azurerm_key_vault_key" "key" {
  name         = "acctestkey%[2]s"
  key_vault_id = azurerm_key_vault.kv.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["encrypt", "decrypt"]
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveKeyVaultKeyStartMoving(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

moved {
  from = azurerm_key_vault_key.key
  to   = azapi_resource.test
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "kv" {
  name                       = "acctestkv%[2]s"
  location                   = azapi_resource.resourceGroup.location
  resource_group_name        = azapi_resource.resourceGroup.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
      "Update",
      "GetRotationPolicy",
      "SetRotationPolicy",
    ]
  }
}

//resource "azurerm_key_vault_key" "key" {
//  name         = "acctestkey%[2]s"
//  key_vault_id = azurerm_key_vault.kv.id
//  key_type     = "RSA"
//  key_size     = 2048
//  key_opts     = ["encrypt", "decrypt"]
//}

resource "azapi_resource" "test" {
  type      = "Microsoft.KeyVault/vaults/keys@2024-11-01"
  parent_id = azurerm_key_vault.kv.id
  name      = "acctestkey%[2]s"
  body = {
    properties = {
      kty     = "RSA"
      keySize = 2048
      keyOps  = ["encrypt", "decrypt"]
    }
  }
  ignore_missing_property = true
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveKeyVaultKeyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

moved {
  from = azurerm_key_vault_key.key
  to   = azapi_resource.test
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "kv" {
  name                       = "acctestkv%[2]s"
  location                   = azapi_resource.resourceGroup.location
  resource_group_name        = azapi_resource.resourceGroup.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
      "Update",
      "GetRotationPolicy",
      "SetRotationPolicy",
    ]
  }
}

//resource "azurerm_key_vault_key" "key" {
//  name         = "acctestkey%[2]s"
//  key_vault_id = azurerm_key_vault.kv.id
//  key_type     = "RSA"
//  key_size     = 2048
//  key_opts     = ["encrypt", "decrypt"]
//}

resource "azapi_resource" "test" {
  type      = "Microsoft.KeyVault/vaults/keys@2024-11-01"
  parent_id = azurerm_key_vault.kv.id
  name      = "acctestkey%[2]s"
  body = {
    properties = {
      kty     = "RSA"
      keySize = 2048
      keyOps  = ["encrypt", "decrypt"]
      rotationPolicy = {
        lifetimeActions = [
          {
            action = {
              type = "rotate"
            }
            trigger = {
              timeAfterCreate = "P90D"
            }
          },

          {
            action = {
              type = "notify"
            }
            trigger = {
              timeBeforeExpiry = "P30D"
            }
          },
        ]
      }
    }
  }
  ignore_missing_property = true
  ignore_casing           = true
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) moveKeyVaultKeyRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

moved {
  from = azurerm_key_vault_key.key
  to   = azapi_resource.test
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "kv" {
  name                       = "acctestkv%[2]s"
  location                   = azapi_resource.resourceGroup.location
  resource_group_name        = azapi_resource.resourceGroup.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
      "Update",
      "GetRotationPolicy",
      "SetRotationPolicy",
    ]
  }
}

//resource "azurerm_key_vault_key" "key" {
//  name         = "acctestkey%[2]s"
//  key_vault_id = azurerm_key_vault.kv.id
//  key_type     = "RSA"
//  key_size     = 2048
//  key_opts     = ["encrypt", "decrypt"]
//}

//resource "azapi_resource" "test" {
//  type      = "Microsoft.KeyVault/vaults/keys@2024-11-01"
//  parent_id = azurerm_key_vault.kv.id
//  name      = "acctestkey%[2]s"
//  body = {
//    properties = {
//      kty     = "RSA"
//      keySize = 2048
//      keyOps  = ["encrypt", "decrypt", "wrapKey"]
//    }
//  }
//  ignore_missing_property = true
//}

removed {
  from = azapi_resource.test
  lifecycle {
    destroy = false
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) SensitiveBody(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

variable "sku_name" {
  type      = string
  default   = "Basic"
  ephemeral = true
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  sensitive_body = {
    properties = {
      sku = {
        name = var.sku_name
      }
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) SensitiveBodyWithHash(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "test" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

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
  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString)
}

func (r GenericResource) SensitiveBodyWithVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "test" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

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

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString)
}

func (r GenericResource) SensitiveBodyWithVersionMultipleTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "test" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

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

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString)
}

func (r GenericResource) SensitiveBodyWithHashMultipleTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "test" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctestRG-%[1]d"
  location = "%[2]s"

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

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString)
}

func (r GenericResource) multipleIdentityIds(data acceptance.TestData, random bool) string {
	identityIds := []string{
		"azapi_resource.userAssignedIdentity1.id",
		"azapi_resource.userAssignedIdentity2.id",
		"azapi_resource.userAssignedIdentity3.id",
	}

	if random {
		// shuffle the identityIds to ensure the order is not fixed
		rand.Shuffle(len(identityIds), func(i, j int) {
			identityIds[i], identityIds[j] = identityIds[j], identityIds[i]
		})
	}
	return fmt.Sprintf(`
%s

resource "azapi_resource" "userAssignedIdentity1" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  name      = "acctest1%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "userAssignedIdentity2" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  name      = "acctest2%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "userAssignedIdentity3" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  name      = "acctest3%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id

  location = "%[3]s"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [%s]
  }

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }

  tags = {
    "Key" = "Value"
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary, strings.Join(identityIds, ", "))
}

func (r GenericResource) automationAccountComplete(data acceptance.TestData, apiVersion string) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@%[3]s"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    identity_ids = []
    type         = "SystemAssigned"
  }
  body = {
    properties = {
      disableLocalAuth = false
      encryption = {
        identity = {
          userAssignedIdentity = null
        }
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        capacity = null
        family   = null
        name     = "Basic"
      }
    }
  }
}
`, r.template(data), data.RandomString, apiVersion)
}

func (r GenericResource) automationAccountCompleteStrictChangeDetection(data acceptance.TestData, apiVersion string) string {
	return fmt.Sprintf(`
%s

provider "azapi" {
  ignore_no_op_changes = false
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@%[3]s"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    identity_ids = []
    type         = "SystemAssigned"
  }
  body = {
    properties = {
      disableLocalAuth = false
      encryption = {
        identity = {
          userAssignedIdentity = null
        }
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        capacity = null
        family   = null
        name     = "Basic"
      }
    }
  }
}
`, r.template(data), data.RandomString, apiVersion)
}

func (r GenericResource) automationAccountCompleteWithNullPropertiesSetup(data acceptance.TestData, apiVersion string) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@%[3]s"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    identity_ids = []
    type         = "SystemAssigned"
  }
  body = {
    properties = {
      disableLocalAuth = false
      encryption = {
        identity = {
          userAssignedIdentity = null
        }
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        capacity = null
        family   = null
        name     = "Basic"
      }
    }
  }
  ignore_null_property = true
}
`, r.template(data), data.RandomString, apiVersion)
}

func (r GenericResource) automationAccountCompleteWithNullProperties(data acceptance.TestData, apiVersion string) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@%[3]s"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    identity_ids = []
    type         = "SystemAssigned"
  }
  body = {
    properties = {
      disableLocalAuth    = false
      encryption          = null
      publicNetworkAccess = null
      sku = {
        capacity = null
        family   = null
        name     = "Basic"
      }
    }
  }
  ignore_null_property = true
}
`, r.template(data), data.RandomString, apiVersion)
}

func (r GenericResource) automationAccountBasic(data acceptance.TestData, apiVersion string) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@%[3]s"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    identity_ids = []
    type         = "SystemAssigned"
  }
  body = {
    properties = {
      disableLocalAuth    = false
      publicNetworkAccess = true
      sku = {
        name = "Basic"
      }
    }
  }
}
`, r.template(data), data.RandomString, apiVersion)
}

func (r GenericResource) automationAccountAzureRM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_automation_account" "automationAccount" {
  name                = "acctest%[2]s"
  location            = azapi_resource.resourceGroup.location
  resource_group_name = azapi_resource.resourceGroup.name
  sku_name            = "Basic"

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) automationAccountAzureRMMovedBasic(data acceptance.TestData, apiVersion string) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

# resource "azurerm_automation_account" "automationAccount" {
#   name      = "acctest%[2]s"
#   location            = azapi_resource.resourceGroup.location
#   resource_group_name = azapi_resource.resourceGroup.name
#   sku_name            = "Basic"
# 
#   identity {
#     type = "SystemAssigned"
#   }
# }

moved {
  from = azurerm_automation_account.automationAccount
  to   = azapi_resource.test
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@%[3]s"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    identity_ids = []
    type         = "SystemAssigned"
  }
  body = {
    properties = {
      disableLocalAuth    = false
      publicNetworkAccess = true
      sku = {
        name = "Basic"
      }
    }
  }
}
`, r.template(data), data.RandomString, apiVersion)
}

func (r GenericResource) automationAccountAzureRMMovedComplete(data acceptance.TestData, apiVersion string) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

# resource "azurerm_automation_account" "automationAccount" {
#   name      = "acctest%[2]s"
#   location            = azapi_resource.resourceGroup.location
#   resource_group_name = azapi_resource.resourceGroup.name
#   sku_name            = "Basic"
# 
#   identity {
#     type = "SystemAssigned"
#   }
# }

moved {
  from = azurerm_automation_account.automationAccount
  to   = azapi_resource.test
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@%[3]s"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    identity_ids = []
    type         = "SystemAssigned"
  }
  body = {
    properties = {
      disableLocalAuth = false
      encryption = {
        identity = {
          userAssignedIdentity = null
        }
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        capacity = null
        family   = null
        name     = "Basic"
      }
    }
  }
}
`, r.template(data), data.RandomString, apiVersion)
}

func (r GenericResource) automationAccountAzureRMMovedCompleteStrictChangeDetection(data acceptance.TestData, apiVersion string) string {
	return fmt.Sprintf(`
%s

provider "azapi" {
  ignore_no_op_changes = false
}

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

# resource "azurerm_automation_account" "automationAccount" {
#   name      = "acctest%[2]s"
#   location            = azapi_resource.resourceGroup.location
#   resource_group_name = azapi_resource.resourceGroup.name
#   sku_name            = "Basic"
# 
#   identity {
#     type = "SystemAssigned"
#   }
# }

moved {
  from = azurerm_automation_account.automationAccount
  to   = azapi_resource.test
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@%[3]s"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  identity {
    identity_ids = []
    type         = "SystemAssigned"
  }
  body = {
    properties = {
      disableLocalAuth = false
      encryption = {
        identity = {
          userAssignedIdentity = null
        }
        keySource = "Microsoft.Automation"
      }
      publicNetworkAccess = true
      sku = {
        capacity = null
        family   = null
        name     = "Basic"
      }
    }
  }
}
`, r.template(data), data.RandomString, apiVersion)
}

func (r GenericResource) modifyPlanSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2019-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/16"
        ]
      }
      subnets = []
    }
  }
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
  schema_validation_enabled = false
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "example-subnet"
  body = {
    properties = {
      addressPrefix                  = "10.0.2.0/24"
      privateEndpointNetworkPolicies = "Disabled"
      defaultOutboundAccess          = false
      delegations                    = []
      routeTable                     = null
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) modifyPlanSubnetUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2019-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/16"
        ]
      }
      subnets = []
    }
  }
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
  schema_validation_enabled = false
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "example-subnet"
  body = {
    properties = {
      addressPrefix                  = "10.0.2.0/24"
      privateEndpointNetworkPolicies = "Disabled"
      defaultOutboundAccess          = false
      serviceEndpoints = [
        {
          service   = "Microsoft.Storage"
          locations = ["westus", "eastus"]
        }
      ]
      delegations = []
      routeTable  = null
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) modifyPlanAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "searchService" {
  type      = "Microsoft.Search/searchServices@2020-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      hostingMode = "default"
      networkRuleSet = {
        ipRules = []
      }
      partitionCount      = 1
      publicNetworkAccess = "enabled"
      replicaCount        = 1
    }
    sku = {
      name = "standard"
    }
  }
  ignore_casing = true
}

data "azapi_resource_action" "getSearchServiceKeys" {
  type        = "Microsoft.Search/searchServices@2023-11-01"
  resource_id = azapi_resource.searchService.id
  action      = "listQueryKeys"
  response_export_values = {
    key = "value[0].key"
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.CognitiveServices/accounts@2025-04-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      apiProperties = {
        qnaAzureSearchEndpointId  = azapi_resource.searchService.id
        qnaAzureSearchEndpointKey = data.azapi_resource_action.getSearchServiceKeys.output.key
      }
    }

    sku = {
      name = "S"
    }
    kind = "TextAnalytics"
  }
  ignore_missing_property = true
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) modifyPlanAccountUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "searchService" {
  type      = "Microsoft.Search/searchServices@2020-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      hostingMode = "default"
      networkRuleSet = {
        ipRules = []
      }
      partitionCount      = 1
      publicNetworkAccess = "enabled"
      replicaCount        = 1
    }
    sku = {
      name = "standard"
    }
  }
  ignore_casing = true
}

data "azapi_resource_action" "getSearchServiceKeys" {
  type        = "Microsoft.Search/searchServices@2023-11-01"
  resource_id = azapi_resource.searchService.id
  action      = "listQueryKeys"
  response_export_values = {
    key = "value[0].key"
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.CognitiveServices/accounts@2025-04-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      apiProperties = {
        qnaAzureSearchEndpointId  = azapi_resource.searchService.id
        qnaAzureSearchEndpointKey = data.azapi_resource_action.getSearchServiceKeys.output.key
        websiteName               = "foo"
      }
    }

    sku = {
      name = "S"
    }
    kind = "TextAnalytics"
  }
  ignore_missing_property = true
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) unknownDiscriminator(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

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

data "azurerm_client_config" "current" {}

resource "azurerm_application_insights" "test" {
  name                = "accappinsights%[2]s"
  location            = azapi_resource.resourceGroup.location
  resource_group_name = azapi_resource.resourceGroup.name
  application_type    = "web"

  lifecycle {
    ignore_changes = [workspace_id]
  }
}

resource "azurerm_key_vault" "test" {
  name                     = "acckeyvault%[2]s"
  location                 = azapi_resource.resourceGroup.location
  resource_group_name      = azapi_resource.resourceGroup.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestsa%[2]s"
  location                        = azapi_resource.resourceGroup.location
  resource_group_name             = azapi_resource.resourceGroup.name
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctestmlws%[2]s"
  location                = azapi_resource.resourceGroup.location
  resource_group_name     = azapi_resource.resourceGroup.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  managed_network {
    isolation_mode = "AllowOnlyApprovedOutbound"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_search_service" "test" {
  name                = "acctestsearch%[2]s"
  location            = azapi_resource.resourceGroup.location
  resource_group_name = azapi_resource.resourceGroup.name
  sku                 = "standard"
}

resource "azapi_resource" "connection" {
  for_each = var.connections

  type      = "Microsoft.MachineLearningServices/workspaces/connections@2025-06-01"
  parent_id = azurerm_machine_learning_workspace.test.id
  name      = each.key
  body = {
    properties = {
      authType      = each.value.type
      category      = "CognitiveSearch"
      expiryTime    = null
      isSharedToAll = true
      metadata = {
        ApiType              = "Azure"
        ApiVersion           = "2024-05-01-preview"
        DeploymentApiVersion = "2023-11-01"
        ResourceId           = azurerm_search_service.test.id
        type                 = "azure_ai_search"
      }
      sharedUserList = []
      target         = "https://${azurerm_search_service.test.name}.search.windows.net/"
    }
  }
}

variable "connections" {
  type = map(object({
    type = string
  }))
  default = {
    "myconnection" = {
      type = "AAD"
    }
  }
}
`, r.template(data), data.RandomString)
}
