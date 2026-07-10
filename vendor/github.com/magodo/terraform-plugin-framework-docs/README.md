# terraform-plugin-framework-docs

A Go library for generating comprehensive documentation for Terraform providers built with [`terraform-plugin-framework`](https://github.com/hashicorp/terraform-plugin-framework).

## Features

- **Complete resource coverage**, including:
    - Provider
    - Managed Resource
    - Data Source
    - Ephemeral Resource
    - List Resource
    - Action
    - Function
- **Enhanced documentation** with descriptions extracted from `CustomType`, `PlanModifiers`, and `Validators`
- **Flexible example integration** — accepts strings that can be derived from acceptance test configurations, ensuring examples remain up-to-date
- **Subcategory support** for organizing generated documents (addresses [#156](https://github.com/hashicorp/terraform-plugin-docs/issues/156))
- **Custom template support** — extensible template system for tailoring documentation output to specific needs
- **`fwdtypes` package** — provides descriptive type wrappers for `terraform-plugin-framework` `basetypes`, enabling documentation for `ObjectType` members (addresses [#333](https://github.com/hashicorp/terraform-plugin-docs/issues/333))

## Comparison with `tfplugindocs`

The official [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) tool relies on Terraform CLI's `terraform providers schema -json` [output](https://developer.hashicorp.com/terraform/cli/commands/providers/schema), which limits the documentation content it can generate. See [this discussion](https://github.com/hashicorp/terraform-plugin-framework/issues/625#issuecomment-1424690927) for details.

This library takes a different approach by reading schemas directly from the provider codebase. To use it, create a separate Go package alongside the provider's `internal` package. See the [example](https://github.com/magodo/terraform-plugin-framework-docs/blob/main/tffwdocs_test.go) for implementation details.

## Example

- [terraform-provider-restful](https://registry.terraform.io/providers/magodo/restful/latest/docs)
