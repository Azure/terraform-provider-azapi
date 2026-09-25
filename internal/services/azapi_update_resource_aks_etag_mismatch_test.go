package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericUpdateResource_aksCustomCATrust(t *testing.T) {
	acceptance.SkipIfCoreAcctestsOnly(t, "AKS clusters require dedicated compute quota")
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			ExternalProviders: aksCustomCATrustExternalProviders(),
			Config:            r.aksCustomCATrust(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func aksCustomCATrustExternalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"azurerm": {
			VersionConstraint: "= 5.0.1",
			Source:            "hashicorp/azurerm",
		},
		"tls": {
			VersionConstraint: "= 4.3.0",
			Source:            "hashicorp/tls",
		},
	}
}

func (r GenericUpdateResource) aksCustomCATrust(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[3]s"
  location            = azapi_resource.resourceGroup.location
  resource_group_name = azapi_resource.resourceGroup.name
  dns_prefix          = "acctestaks%[3]s"

  default_node_pool {
    name       = "system"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
  }

  node_provisioning_profile {
    mode = "Manual"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "tls_private_key" "test" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "test" {
  private_key_pem       = tls_private_key.test.private_key_pem
  validity_period_hours = 8760
  is_ca_certificate     = true

  subject {
    common_name = "acctestaks%[3]s"
  }

  allowed_uses = [
    "cert_signing",
    "crl_signing",
  ]
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.ContainerService/managedClusters@2024-10-02-preview"
  resource_id = azurerm_kubernetes_cluster.test.id
  locks       = [azurerm_kubernetes_cluster.test.id]

  body = {
    properties = {
      agentPoolProfiles = [{
        name                = "system"
        enableCustomCATrust = true
      }]
      securityProfile = {
        customCATrustCertificates = [
          base64encode(tls_self_signed_cert.test.cert_pem),
        ]
      }
    }
  }
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString)
}
