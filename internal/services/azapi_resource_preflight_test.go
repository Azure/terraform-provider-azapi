package services_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericResource_preflightMockPropertyValue(t *testing.T) {
	// in this test, the parent_id, name and a field inside the properties bag are unknown,
	// these unknown values are replaced with valid mock values to allow the preflight validation to pass
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.preflightMockPropertyValue(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_preflightUnknownValueOutsidePropertiesBag(t *testing.T) {
	// in this test, there's an unknown value outside the properties bag,
	// preflight validation should skip this resource
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.preflightMockPropertyValueOutsidePropertiesBag(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_preflightValidationFailed(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.preflightFeatureFlag(data, true),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("AccountTypeMissing"),
		},
	})
}

func TestAccGenericResource_preflightDisableValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.preflightFeatureFlag(data, true),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("AccountTypeMissing"),
		},
		{
			Config:             r.preflightFeatureFlag(data, false),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_preflightExtensionResourceValidation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.preflightExtensionResourceValidation(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_preflightWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.preflightWithIdentity(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_preflightTenantScopedResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.preflightTenantScopedResource(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_preflightSubscriptionScopedResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.preflightSubscriptionScopedResource(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_preflightManagementGroupScopedNestedResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.preflightManagementGroupScopedNestedResource(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_preflightResourceGroupScopedTrackedResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.preflightResourceGroupScopedTrackedResource(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccGenericResource_preflightResourceGroupScopedNestedResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.preflightResourceGroupScopedNestedResource(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	})
}

func (r GenericResource) preflightMockPropertyValue(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azapi" {
  enable_preflight = true
}

%[1]s

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = azapi_resource.resourceGroup.id
  location  = "westus"
  body = {
    kind = "StorageV2"
    properties = {
      accessTier = azapi_resource.resourceGroup.id
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) preflightMockPropertyValueOutsidePropertiesBag(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azapi" {
  enable_preflight = true
}

%[1]s

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = azapi_resource.resourceGroup.id
  location  = "westus"
  body = {
    kind = azapi_resource.resourceGroup.id
    properties = {
      accessTier = "Cold"
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) preflightFeatureFlag(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azapi" {
  enable_preflight = %[2]t
}

%[1]s

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = azapi_resource.resourceGroup.id
  location  = "westus"
  body = {
    kind = "StorageV2"
  }
  schema_validation_enabled = false
}
`, r.template(data), enabled)
}

func (r GenericResource) preflightExtensionResourceValidation(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azapi" {
  enable_preflight = true
}

%s
`, r.extensionScope(data))
}

func (r GenericResource) preflightWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azapi" {
  enable_preflight = true
}

%[1]s

resource "azapi_resource" "aksCluster" {
  type      = "Microsoft.ContainerService/managedClusters@2024-06-02-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = azapi_resource.resourceGroup.id
  location  = "westus"
  identity {
    type = "SystemAssigned"
  }
  body = {
    properties = {
      agentPoolProfiles = [
        {
          count  = 1
          mode   = "System"
          name   = "default"
          vmSize = "Standard_DS2_v2"
        },
      ]
      dnsPrefix = "exampleaks"
    }
  }
  schema_validation_enabled = false
}
`, r.template(data))
}

func (r GenericResource) preflightTenantScopedResource(_ acceptance.TestData) string {
	return `
provider "azapi" {
  enable_preflight = true
}
resource "azapi_resource" "billingAccount" {
  type      = "Microsoft.Billing/billingAccounts@2024-04-01"
  name      = "PCN.1234567890"
  parent_id = "/"
  body = {

  }
}

resource "azapi_resource" "managementGroup" {
  type      = "Microsoft.Management/managementGroups@2023-04-01"
  name      = "myManagementGroup"
  parent_id = "/"
  body = {
    properties = {
      displayName = "My Management Group"
    }
  }
}
`
}

func (r GenericResource) preflightSubscriptionScopedResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azapi" {
  enable_preflight = true
}

%[1]s
`, r.template(data))
}

func (r GenericResource) preflightResourceGroupScopedTrackedResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azapi" {
  enable_preflight = true
}

%[1]s

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2021-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "myVirtualNetwork"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}
`, r.template(data))
}

func (r GenericResource) preflightResourceGroupScopedNestedResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azapi" {
  enable_preflight = true
}

%[1]s

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2021-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "myVirtualNetwork"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2021-02-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "mySubnet"
  body = {
    properties = {
      addressPrefix = "10.0.0.0/24"
    }
  }
}
`, r.template(data))
}

func (r GenericResource) preflightManagementGroupScopedNestedResource(_ acceptance.TestData) string {
	return `
provider "azapi" {
  enable_preflight = true
}

resource "azapi_resource" "billingAccount" {
  type      = "Microsoft.Billing/billingAccounts@2024-04-01"
  name      = "PCN.1234567"
  parent_id = "/providers/Microsoft.Management/managementGroups/myManagementGroup"
  body = {
  }
}
`
}
