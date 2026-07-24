package migration

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func telemetryHeadersAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"module_id":      types.StringType,
		"module_version": types.StringType,
	}
}
