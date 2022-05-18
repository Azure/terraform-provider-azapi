## Unreleased
FEATURES:

ENHANCEMENTS:

BUG FIXES:
- Improve error message for schema validation failure.
- DefaultAzureCredential reads the client ID of a user-assigned managed identity.
- Fix the modification is not working, when use `azapi_update_resource` to modify additional properties.

## v0.2.0
FEATURES:

ENHANCEMENTS:
- Setting `response_export_values = ["*"]` will export the full response body.

BUG FIXES:
- Fix incorrect ID format in the imported `azapi_resource` resource. 
- Fix incorrect `body` content in the imported `azapi_resource` resource.

## v0.1.1

FEATURES:

ENHANCEMENTS:

BUG FIXES:

- Fix document format.

## v1.1.0

FEATURES:
- **New Data Source**: azapi_resource
- **New Resource**: azapi_resource
- **New Resource**: azapi_update_resource
- **Provider feature**: support default location and default tags

ENHANCEMENTS:

BUG FIXES:
