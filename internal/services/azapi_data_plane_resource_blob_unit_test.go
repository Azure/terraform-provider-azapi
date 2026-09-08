package services

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const blobResourceType = "Microsoft.Storage/storageAccounts/blobServices/containers/blobs@2023-11-03"

func TestValidateDataPlaneResourceBodyModeForBlob(t *testing.T) {
	body := dynamicObject(map[string]attr.Value{
		"content_base64": types.StringValue("YmxvYg=="),
	})
	version := types.MapValueMust(types.StringType, map[string]attr.Value{
		"payload": types.StringValue("v1"),
	})

	tests := []struct {
		name          string
		body          types.Dynamic
		sensitiveBody types.Dynamic
		version       types.Map
		wantError     string
	}{
		{
			name:          "regular body",
			body:          body,
			sensitiveBody: types.DynamicNull(),
			version:       types.MapNull(types.StringType),
		},
		{
			name:          "sensitive body with version",
			body:          emptyDynamicObject(),
			sensitiveBody: body,
			version:       version,
		},
		{
			name:          "mixed bodies",
			body:          body,
			sensitiveBody: body,
			version:       version,
			wantError:     "requires exactly one of body or sensitive_body",
		},
		{
			name:          "sensitive body without version",
			body:          emptyDynamicObject(),
			sensitiveBody: body,
			version:       types.MapNull(types.StringType),
			wantError:     "sensitive_body_version must be set",
		},
		{
			name:          "missing payload",
			body:          emptyDynamicObject(),
			sensitiveBody: types.DynamicNull(),
			version:       types.MapNull(types.StringType),
			wantError:     "requires exactly one of body or sensitive_body",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateDataPlaneResourceBodyMode(&DataPlaneResourceModel{
				Type:                 types.StringValue(blobResourceType),
				Body:                 test.body,
				SensitiveBody:        test.sensitiveBody,
				SensitiveBodyVersion: test.version,
			})
			if test.wantError == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("expected error containing %q, got %v", test.wantError, err)
			}
		})
	}
}

func TestShouldSkipExclusiveSensitiveBodyUpdate(t *testing.T) {
	version := types.MapValueMust(types.StringType, map[string]attr.Value{
		"payload": types.StringValue("v1"),
	})
	state := &DataPlaneResourceModel{
		Type:                 types.StringValue(blobResourceType),
		SensitiveBodyVersion: version,
	}
	plan := &DataPlaneResourceModel{
		Type:                 types.StringValue(blobResourceType),
		SensitiveBodyVersion: version,
	}
	if !shouldSkipExclusiveSensitiveBodyUpdate(plan, state) {
		t.Fatal("expected unchanged sensitive body version to skip the upload")
	}

	plan.SensitiveBodyVersion = types.MapValueMust(types.StringType, map[string]attr.Value{
		"payload": types.StringValue("v2"),
	})
	if shouldSkipExclusiveSensitiveBodyUpdate(plan, state) {
		t.Fatal("expected changed sensitive body version to upload")
	}
}

func dynamicObject(values map[string]attr.Value) types.Dynamic {
	typesByName := make(map[string]attr.Type, len(values))
	for name, value := range values {
		typesByName[name] = value.Type(context.Background())
	}
	return types.DynamicValue(types.ObjectValueMust(typesByName, values))
}

func emptyDynamicObject() types.Dynamic {
	return types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}))
}
