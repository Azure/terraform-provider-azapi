package services

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestWithTelemetryHeaders(t *testing.T) {
	headers := map[string]string{
		"X-Custom-Header":    "custom",
		"x-module-id":        "from-headers-lower",
		"X-MODULE-VERSION":   "from-headers-upper",
		"x-module-version":   "from-headers-lower-version",
		"another-header-key": "another-value",
	}
	telemetryHeaders := types.ObjectValueMust(telemetryHeadersAttributeTypes(), map[string]attr.Value{
		"module_id":      types.StringValue("Azure/avm-res-network-virtualnetwork/azurerm"),
		"module_version": types.StringValue("0.19.0"),
	})

	merged := withTelemetryHeaders(headers, telemetryHeaders)
	if got := merged["X-Custom-Header"]; got != "custom" {
		t.Fatalf("expected custom header to be preserved, got %q", got)
	}
	if got := merged[telemetryHeaderModuleID]; got != "Azure/avm-res-network-virtualnetwork/azurerm" {
		t.Fatalf("expected %s header to be set, got %q", telemetryHeaderModuleID, got)
	}
	if got := merged[telemetryHeaderModuleVersion]; got != "0.19.0" {
		t.Fatalf("expected %s header to be set, got %q", telemetryHeaderModuleVersion, got)
	}
	if _, ok := merged["x-module-id"]; ok {
		t.Fatalf("expected lowercase module id header key to be removed")
	}
	if _, ok := merged["x-module-version"]; ok {
		t.Fatalf("expected lowercase module version header key to be removed")
	}
	if _, ok := headers[telemetryHeaderModuleID]; ok {
		t.Fatalf("expected source header map to remain unchanged")
	}
}
