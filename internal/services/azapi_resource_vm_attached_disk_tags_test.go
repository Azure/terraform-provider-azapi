package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericResource_vmWithAttachedDiskTagsRemoval(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.vmWithAttachedDiskTags(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output.tags.environment").HasValue("acctest"),
				check.That(data.ResourceName).Key("output.tags.removable").HasValue("value"),
			),
		},
		{
			Config: r.vmWithAttachedDiskTags(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.tags.removable").DoesNotExist(),
				check.That(data.ResourceName).Key("output.tags.environment").HasValue("acctest"),
			),
		},
		{
			Config:   r.vmWithAttachedDiskTags(data, false),
			PlanOnly: true,
		},
	})
}

func (r GenericResource) vmWithAttachedDiskTags(data acceptance.TestData, includeRemovableTag bool) string {
	removableTag := ""
	if includeRemovableTag {
		removableTag = `removable = "value"`
	}

	return fmt.Sprintf(`
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2024-03-01"
  name     = "acctestrg%s"
  location = %q
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  name      = "acctest-vnet-%s"
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

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  name      = "acctest-subnet"
  parent_id = azapi_resource.virtualNetwork.id

  body = {
    properties = {
      addressPrefix = "10.0.1.0/24"
    }
  }
}

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2024-05-01"
  name      = "acctest-nic-%s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      ipConfigurations = [{
        name = "internal"
        properties = {
          primary                   = true
          privateIPAddressVersion   = "IPv4"
          privateIPAllocationMethod = "Dynamic"
          subnet = {
            id = azapi_resource.subnet.id
          }
        }
      }]
    }
  }
}

resource "azapi_resource" "attachedDisk" {
  type      = "Microsoft.Compute/disks@2023-04-02"
  name      = "acctest-disk-%s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      creationData = {
        createOption = "Empty"
      }
      diskSizeGB = 4
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Compute/virtualMachines@2023-03-01"
  name      = "acctest-vm-%s"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location

  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_D2s_v3"
      }
      networkProfile = {
        networkInterfaces = [{
          id = azapi_resource.networkInterface.id
          properties = {
            primary = true
          }
        }]
      }
      osProfile = {
        adminPassword = "Acctest-Passw0rd!"
        adminUsername = "acctest"
        computerName  = "acctest%s"
        linuxConfiguration = {
          disablePasswordAuthentication = false
        }
      }
      storageProfile = {
        imageReference = {
          offer     = "0001-com-ubuntu-server-jammy"
          publisher = "Canonical"
          sku       = "22_04-lts-gen2"
          version   = "latest"
        }
        osDisk = {
          caching      = "ReadWrite"
          createOption = "FromImage"
          managedDisk = {
            storageAccountType = "Standard_LRS"
          }
        }
        dataDisks = [{
          caching      = "ReadWrite"
          createOption = "Attach"
          lun          = 0
          managedDisk = {
            id = azapi_resource.attachedDisk.id
          }
        }]
      }
    }
    tags = {
      environment = "acctest"
      %s
    }
  }

  response_export_values = ["tags"]
}
`,
		data.RandomString,
		data.LocationPrimary,
		data.RandomString,
		data.RandomString,
		data.RandomString,
		data.RandomString,
		data.RandomString,
		removableTag,
	)
}
