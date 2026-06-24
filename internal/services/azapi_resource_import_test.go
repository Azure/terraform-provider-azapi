package services_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericResource_importIdWithApiVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azapi_resource.source").ExistsInAzure(r),
			),
		},
		{
			Config: r.importIdWithApiVersion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azapi_resource.import_id_with_api_version").ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_importIdWithoutApiVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azapi_resource.source").ExistsInAzure(r),
			),
		},
		{
			Config:      r.importIdWithoutApiVersion(data),
			ExpectError: regexp.MustCompile("missing the `api-version` query parameter"),
		},
	})
}

func TestAccGenericResource_importIdAndType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azapi_resource.source").ExistsInAzure(r),
			),
		},
		{
			Config: r.importIdAndType(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azapi_resource.import_id_and_type").ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_importIdClassicWithApiVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azapi_resource.source").ExistsInAzure(r),
			),
		},
		{
			Config: r.importIdClassicWithApiVersion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azapi_resource.import_id_classic").ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericResource_importIdClassicWithoutApiVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azapi_resource.source").ExistsInAzure(r),
			),
		},
		{
			Config:      r.importIdClassicWithoutApiVersion(data),
			ExpectError: regexp.MustCompile("missing the `api-version` query parameter"),
		},
	})
}

func (r GenericResource) importBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2024-03-01"
  name     = "acctestrg%[1]s"
  location = "eastus"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}

resource "azapi_resource" "source" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctestsource%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}
`, data.RandomString)
}

// Case 1a: Identity-based import with only ID (ID contains API version as query parameter)
func (r GenericResource) importIdWithApiVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2024-03-01"
  name     = "acctestrg%[1]s"
  location = "eastus"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}

removed {
  from = azapi_resource.source
  lifecycle {
    destroy = false
  }
}

locals {
  source_id = "${azapi_resource.resourceGroup.id}/providers/Microsoft.Network/virtualNetworks/acctestsource%[1]s"
}

import {
  to = azapi_resource.import_id_with_api_version
  identity = {
    id = "${local.source_id}?api-version=2024-05-01"
  }
}

resource "azapi_resource" "import_id_with_api_version" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctestsource%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}
`, data.RandomString)
}

// Case 1b: Identity-based import with only ID (ID does NOT contain API version)
func (r GenericResource) importIdWithoutApiVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2024-03-01"
  name     = "acctestrg%[1]s"
  location = "eastus"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}

removed {
  from = azapi_resource.source
  lifecycle {
    destroy = false
  }
}

locals {
  source_id = "${azapi_resource.resourceGroup.id}/providers/Microsoft.Network/virtualNetworks/acctestsource%[1]s"
}

import {
  to = azapi_resource.import_id_without_api_version
  identity = {
    id = local.source_id
  }
}

resource "azapi_resource" "import_id_without_api_version" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctestsource%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}
`, data.RandomString)
}

// Case 2: Identity-based import with both ID and Type
func (r GenericResource) importIdAndType(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2024-03-01"
  name     = "acctestrg%[1]s"
  location = "eastus"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}

removed {
  from = azapi_resource.source
  lifecycle {
    destroy = false
  }
}

locals {
  source_id = "${azapi_resource.resourceGroup.id}/providers/Microsoft.Network/virtualNetworks/acctestsource%[1]s"
}

import {
  to = azapi_resource.import_id_and_type
  identity = {
    id   = local.source_id
    type = "Microsoft.Network/virtualNetworks@2024-05-01"
  }
}

resource "azapi_resource" "import_id_and_type" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctestsource%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}
`, data.RandomString)
}

// Case 3a: Classic import with only ID (ID contains API version as query parameter)
func (r GenericResource) importIdClassicWithApiVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2024-03-01"
  name     = "acctestrg%[1]s"
  location = "eastus"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}

removed {
  from = azapi_resource.source
  lifecycle {
    destroy = false
  }
}

locals {
  source_id = "${azapi_resource.resourceGroup.id}/providers/Microsoft.Network/virtualNetworks/acctestsource%[1]s"
}

import {
  to = azapi_resource.import_id_classic
  id = "${local.source_id}?api-version=2024-05-01"
}

resource "azapi_resource" "import_id_classic" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctestsource%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}
`, data.RandomString)
}

// Case 3b: Classic import with only ID (ID does NOT contain API version)
func (r GenericResource) importIdClassicWithoutApiVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2024-03-01"
  name     = "acctestrg%[1]s"
  location = "eastus"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}

removed {
  from = azapi_resource.source
  lifecycle {
    destroy = false
  }
}

locals {
  source_id = "${azapi_resource.resourceGroup.id}/providers/Microsoft.Network/virtualNetworks/acctestsource%[1]s"
}

import {
  to = azapi_resource.import_id_classic
  id = local.source_id
}

resource "azapi_resource" "import_id_classic" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctestsource%[1]s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
}
`, data.RandomString)
}
