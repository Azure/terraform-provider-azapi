package services

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const telemetryHeadersMarkdownDescription = "Optional telemetry headers to include in requests. Values are sent as request headers but not persisted in Terraform state. Supported fields: `module_id` (`X-MODULE-ID`) and `module_version` (`X-MODULE-VERSION`). For stateful resources, these headers are applied on create and update requests."

const (
	telemetryHeaderModuleID      = "X-MODULE-ID"
	telemetryHeaderModuleVersion = "X-MODULE-VERSION"
)

func telemetryHeadersAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"module_id":      types.StringType,
		"module_version": types.StringType,
	}
}

func withTelemetryHeaders(headers map[string]string, telemetryHeaders types.Object) map[string]string {
	merged := make(map[string]string, len(headers)+2)
	for k, v := range headers {
		merged[k] = v
	}

	if telemetryHeaders.IsNull() || telemetryHeaders.IsUnknown() {
		return merged
	}

	attributes := telemetryHeaders.Attributes()
	moduleID, hasModuleID := attributes["module_id"].(types.String)
	moduleVersion, hasModuleVersion := attributes["module_version"].(types.String)
	hasModuleID = hasModuleID && !moduleID.IsNull() && !moduleID.IsUnknown()
	hasModuleVersion = hasModuleVersion && !moduleVersion.IsNull() && !moduleVersion.IsUnknown()

	if hasModuleID {
		for k := range merged {
			if strings.EqualFold(k, telemetryHeaderModuleID) {
				delete(merged, k)
			}
		}
		merged[telemetryHeaderModuleID] = moduleID.ValueString()
	}
	if hasModuleVersion {
		for k := range merged {
			if strings.EqualFold(k, telemetryHeaderModuleVersion) {
				delete(merged, k)
			}
		}
		merged[telemetryHeaderModuleVersion] = moduleVersion.ValueString()
	}
	return merged
}
