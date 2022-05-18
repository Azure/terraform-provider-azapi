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

func (r GenericUpdateResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resourceType := state.Attributes["type"]
	id, err := parse.NewResourceID(state.ID, resourceType)
	if err != nil {
		return nil, err
	}

	resp, err := client.ResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			exist := false
			return &exist, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	exist := len(utils.GetId(resp)) != 0
	return &exist, nil
}

func (r GenericUpdateResource) automationAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id = azurerm_automation_account.test.id
  body = jsonencode({
    properties = {
      publicNetworkAccess = true
    }
  })
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (r GenericUpdateResource) automationAccountWithNameParentId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_update_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  name      = azurerm_automation_account.test.name
  parent_id = azurerm_resource_group.test.id
  body = jsonencode({
    properties = {
      publicNetworkAccess = true
    }
  })
}
`, r.template(data), data.RandomStringOfLength(5))
}

func (GenericUpdateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    azurerm = {
      version = "= 2.75.0"
      source  = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.LocationPrimary, data.RandomStringOfLength(10))
}
