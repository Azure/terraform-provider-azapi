## Unreleased
FEATURES:

ENHANCEMENTS:
- `azapi_resource`: Supports `locks` which used to prevent creating/modifying/deleting resources at the same time.
- `azapi_update_resource`: Supports `locks` which used to prevent creating/modifying/deleting resources at the same time.

BUG FIXES:

## v0.3.0
FEATURES:

ENHANCEMENTS:
- Update bicep types to https://github.com/Azure/bicep-types-az/commit/644ff521c92ce8d493f6da977af12377f32abffc

BUG FIXES:

## v0.2.1
FEATURES:

ENHANCEMENTS:

BUG FIXES:
- Improve error message for schema validation failure.
- DefaultAzureCredential reads the client ID of a user-assigned managed identity.
- Fix the modification is not working, when use `azapi_update_resource` to modify additional properties.
- Fix crash when use `azapi_update_resource` on a resource whose id is null
- Fix crash when the discriminated type is not in the embedded schema

## v0.2.0
FEATURES:

ENHANCEMENTS:
- Setting `response_export_values = ["*"]` will export the full response body.
- Update bicep types to https://github.com/Azure/bicep-types-az/commit/57f3ecc750648562cf170ef456ef39533872b101

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
