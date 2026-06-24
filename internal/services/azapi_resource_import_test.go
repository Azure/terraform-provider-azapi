package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericResource_importWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.importWithIdentityAllCases(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azapi_resource.import_id_with_api_version").ExistsInAzure(r),
				check.That("azapi_resource.import_id_without_api_version").ExistsInAzure(r),
				check.That("azapi_resource.import_id_and_type").ExistsInAzure(r),
			),
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

resource "azapi_resource" "source_a" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctesta%[1]s"
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

resource "azapi_resource" "source_b" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctestb%[1]s"
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

resource "azapi_resource" "source_c" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctestc%[1]s"
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

func (r GenericResource) importWithIdentityAllCases(data acceptance.TestData) string {
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

# The source_a/b/c resources were created in the previous step. Forget them from state
# without destroying the Azure resources so each can be re-adopted via identity-based
# import into a dedicated address below. This ensures every physical resource is managed
# by exactly one Terraform address, avoiding concurrent deletes during destroy.
removed {
  from = azapi_resource.source_a
  lifecycle {
    destroy = false
  }
}

removed {
  from = azapi_resource.source_b
  lifecycle {
    destroy = false
  }
}

removed {
  from = azapi_resource.source_c
  lifecycle {
    destroy = false
  }
}

locals {
  source_a_id = "${azapi_resource.resourceGroup.id}/providers/Microsoft.Network/virtualNetworks/acctesta%[1]s"
  source_b_id = "${azapi_resource.resourceGroup.id}/providers/Microsoft.Network/virtualNetworks/acctestb%[1]s"
  source_c_id = "${azapi_resource.resourceGroup.id}/providers/Microsoft.Network/virtualNetworks/acctestc%[1]s"
}

# Case 1a: Identity-based import with only ID (ID contains API version as query parameter)
import {
  to = azapi_resource.import_id_with_api_version
  identity = {
    id = "${local.source_a_id}?api-version=2024-05-01"
  }
}

resource "azapi_resource" "import_id_with_api_version" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctesta%[1]s"
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

# Case 1b: Identity-based import with only ID (ID does NOT contain API version)
import {
  to = azapi_resource.import_id_without_api_version
  identity = {
    id = local.source_b_id
  }
}

resource "azapi_resource" "import_id_without_api_version" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctestb%[1]s"
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

# Case 2: Identity-based import with both ID and Type
import {
  to = azapi_resource.import_id_and_type
  identity = {
    id   = local.source_c_id
    type = "Microsoft.Network/virtualNetworks@2024-05-01"
  }
}

resource "azapi_resource" "import_id_and_type" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctestc%[1]s"
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
