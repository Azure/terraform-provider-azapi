# Contributing to the AzAPI Provider

Thank you for your interest in contributing to the AzAPI Terraform provider! This guide will help you understand how to contribute data-plane resource support to the provider.

## Contributing Data-Plane Resource Support

The AzAPI provider supports Azure data-plane resources through the `azapi_data_plane_resource` resource type. This section describes how to add support for new data-plane resources or customize behavior for existing ones.

### Overview

Data-plane resources are Azure resources that are managed through data-plane APIs rather than Azure Resource Manager (ARM) APIs. Examples include:
- Azure App Configuration key-values
- Azure Key Vault secrets and keys
- Azure Synapse workspace datasets
- Azure Purview collections
- Azure IoT Central applications

### Architecture

The data-plane resource support in AzAPI consists of three main components:

1. **Resource Type Registration** - Mapping resource types to their API URL formats
2. **Data Plane Client** - HTTP client for making data-plane API calls
3. **Resource Customization** (Optional) - Custom CRUD operations for resources that don't follow standard patterns

### Step 1: Register the Resource Type

Data-plane resource types must be registered in the `data_plane_resources.json` file to define their URL format and parent ID structure.

**Location:** `internal/services/parse/data_plane_resources.json`

Add a new entry to the JSON array:

```json
{
  "UrlFormat": "{parentId}/path/to/{name}",
  "ResourceType": "Microsoft.ServiceName/parentResourceType/resourceType",
  "ParentIDExample": "{serviceName}.serviceEndpoint.azure.net",
  "Url": "/path/to/{resourceName}"
}
```

**Field Descriptions:**

- **`UrlFormat`**: Template for constructing the full resource URL
  - Use `{parentId}` as a placeholder for the parent resource endpoint
  - Use `{name}` as a placeholder for the resource name
  - Example: `{parentId}/kv/{name}` for App Configuration key-values

- **`ResourceType`**: The full Azure resource type
  - Format: `Microsoft.ServiceName/parentResourceType/childResourceType`
  - Must match the type users will specify in their Terraform configurations
  - Example: `Microsoft.AppConfiguration/configurationStores/keyValues`

- **`ParentIDExample`**: Documentation showing the expected parent ID format
  - Helps users understand what to pass as the `parent_id` parameter
  - Example: `{storeName}.azconfig.io`

- **`Url`**: The actual API path (for documentation purposes)
  - Shows the final URL structure
  - Example: `/kv/{key}`

### Step 2: Test the Resource

After registering the resource type, you can immediately use it with `azapi_data_plane_resource` in your Terraform configurations:

```hcl
resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.ServiceName/parentResourceType/resourceType@2023-01-01"
  parent_id = "{serviceName}.serviceEndpoint.azure.net"
  name      = "example-resource"
  
  body = {
    properties = {
      # Your resource properties here
    }
  }
}
```

The AzAPI provider will automatically:
- Construct the proper URL using the registered `UrlFormat`
- Make HTTP PUT requests for create/update operations
- Make HTTP GET requests for read operations
- Make HTTP DELETE requests for delete operations
- Handle authentication using Azure credentials

### Step 3: Add Customization (If Needed)

Most data-plane resources work with the default CRUD operations (PUT for create/update, GET for read, DELETE for delete). However, some resources require custom behavior.

**When to Add Customization:**

- The API uses non-standard HTTP methods (e.g., POST for create)
- The API uses different endpoints for different operations
- The API requires special action paths (e.g., `/create`, `/rotate`)
- The API has unique authentication or header requirements

**How to Add Customization:**

#### 3.1 Create a Customization File

Create a new file in `internal/services/customization/` following the naming pattern `{service}_customization.go`.

**Location:** `internal/services/customization/{resource_type}_customization.go`

Example: `internal/services/customization/key_vault_key_customization.go`

#### 3.2 Implement the DataPlaneResource Interface

```go
package customization

import (
	"context"
	"net/http"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

type MyResourceCustomization struct {
}

// GetResourceType returns the resource type this customization applies to
func (r MyResourceCustomization) GetResourceType() string {
	return "Microsoft.ServiceName/parentResourceType/resourceType"
}

// CreateFunc returns a custom create function
// Return nil to use the default behavior
func (r MyResourceCustomization) CreateFunc() CreateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		// Custom create logic
		// Example: Use POST with /create action instead of PUT
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "create", id.ApiVersion, http.MethodPost, body, options)
		return err
	}
}

// UpdateFunc returns a custom update function
// Return nil to use the default behavior
func (r MyResourceCustomization) UpdateFunc() UpdateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		// Custom update logic
		// Example: Use PATCH instead of PUT
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "", id.ApiVersion, http.MethodPatch, body, options)
		return err
	}
}

// ReadFunc returns a custom read function
// Return nil to use the default behavior
func (r MyResourceCustomization) ReadFunc() ReadFunc {
	// Return nil if the default GET behavior is sufficient
	return nil
}

// DeleteFunc returns a custom delete function
// Return nil to use the default behavior
func (r MyResourceCustomization) DeleteFunc() DeleteFunc {
	// Return nil if the default DELETE behavior is sufficient
	return nil
}

// Ensure the struct implements the interface
var _ DataPlaneResource = &MyResourceCustomization{}
```

#### 3.3 Register the Customization

Add your customization to the registry in `internal/services/customization/registration.go`:

```go
func init() {
	// Existing registrations...
	
	// Register your customization
	var myResourceCustomization DataPlaneResource = MyResourceCustomization{}
	customizations[strings.ToLower(myResourceCustomization.GetResourceType())] = myResourceCustomization
}
```

### Step 4: Add Tests

Create acceptance tests for your data-plane resource:

- **Without customization**: Add tests to `internal/services/azapi_data_plane_resource_test.go`
- **With customization**: Add tests to `internal/services/azapi_data_plane_resource_customization_test.go`

**Test Structure:**

```go
func TestAccDataPlaneResource_myResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.myResourceBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DataPlaneResource) myResourceBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azapi" {
}

# Create parent resource if needed
resource "azapi_resource" "parent" {
  type      = "Microsoft.ServiceName/parentResourceType@2023-01-01"
  parent_id = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}"
  name      = "acctest%[2]s"
  location  = "%[1]s"
  body = {
    properties = {}
  }
}

# Create data-plane resource
resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.ServiceName/parentResourceType/resourceType@2023-01-01"
  parent_id = trimprefix(azapi_resource.parent.output.properties.dataEndpoint, "https://")
  name      = "acctest-resource"
  
  body = {
    properties = {
      key = "value"
    }
  }
}
`, data.LocationPrimary, data.RandomString)
}
```

### Step 5: Add Documentation Examples

Add example configurations to help users understand how to use the resource.

**Location:** `examples/{ResourceType}`

Create a directory with the resource type and API version, then add example files:

```
examples/
└── Microsoft.ServiceName_parentResourceType_resourceType/
    └── main.tf           # Basic usage example
```

**Example `main.tf`:**

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

# Parent resource
resource "azapi_resource" "parent" {
  type      = "Microsoft.ServiceName/parentResourceType@2023-01-01"
  parent_id = "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}"
  name      = "example-parent"
  location  = "eastus"
  body = {
    properties = {}
  }
}

# Data-plane resource
resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.ServiceName/parentResourceType/resourceType@2023-01-01"
  parent_id = trimprefix(azapi_resource.parent.output.properties.dataEndpoint, "https://")
  name      = "example-resource"
  
  body = {
    properties = {
      description = "Example resource"
      enabled     = true
    }
  }
}

output "resource_id" {
  value = azapi_data_plane_resource.example.id
}
```

After adding examples, run `make docs` to regenerate the documentation:

```bash
make docs
```

This will update the documentation in the `docs/` directory based on your examples and code changes.


