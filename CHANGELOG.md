## v1.2.0 (Unreleased)
FEATURES:
- `azapi`: Support `client_certificate_password` option.

ENHANCEMENTS:

BUG FIXES:

## v1.1.0
FEATURES:

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/e641570bedc5004498d3e374adb60fdfd3521b09

BUG FIXES:
- `azapi_resource_action`: The `output` is not refreshed when `body` is changed.

## v1.0.0
FEATURES:

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/a6dabb0cd645c17a1accf3ec1be4d7930e982b23

BUG FIXES:

BREAKING CHANGES:
- `azapi_resource`: `ignore_missing_property`'s default value changed from `false` to `true`.
- `azapi_update_resource`: `ignore_missing_property`'s default value changed from `false` to `true`.

## v0.6.0
FEATURES:

ENHANCEMENTS:
- `azapi_resource_action`: Supports `locks` which used to prevent modifying resources at the same time.
- `azapi_resource_action`: Supports parse response which `Content-Type` is `text/plain`.
- Improve validation on `type`, `parent_id` and `resource_id`.
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/5451fcd5e1bf4d8313d475d8e3dc28efc7a77e2a

BUG FIXES:

## v0.5.0
FEATURES:
- **New Data Source**: azapi_resource_action
- **New Resource**: azapi_resource_action

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/813d8bbc9ecf432a2a0ff2769627592fae34369f

BUG FIXES:
- DefaultAzureCredential authentication failed because empty clientId is set

## v0.4.0
FEATURES:

ENHANCEMENTS:
- `azapi_resource`: Supports default api-version when importing existing resource into terraform state.
- `azapi_resource`: Supports `locks` which used to prevent creating/modifying/deleting resources at the same time.
- `azapi_update_resource`: Supports `locks` which used to prevent creating/modifying/deleting resources at the same time.
- `azapi_resource` data source: Supports configuring `resource_id` as an alternative way to configure `name` and `parent_id`.
- `azapi` provider: Supports `partner_id`, `disable_terraform_partner_id` and `disable_terraform_partner_id`.
- Update bicep types to https://github.com/Azure/bicep-types-az/commit/ea703e2aba0d1c024f33124ee2cd34bc0c6084b5

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
