package services

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	fwaction "github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	fwephemeral "github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephemeralschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAzapiResourceActionSchemaSensitiveBody(t *testing.T) {
	var response fwaction.SchemaResponse
	a := AzapiResourceAction{}
	a.Schema(context.Background(), fwaction.SchemaRequest{}, &response)

	sensitiveBody, ok := response.Schema.Attributes["sensitive_body"].(actionschema.DynamicAttribute)
	if !ok {
		t.Fatalf("expected sensitive_body to be a dynamic attribute")
	}
	if !sensitiveBody.WriteOnly {
		t.Fatalf("expected sensitive_body to be write-only")
	}
}

func TestActionResourceSchemaSensitiveBody(t *testing.T) {
	var response fwresource.SchemaResponse
	r := ActionResource{}
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &response)

	sensitiveBody, ok := response.Schema.Attributes["sensitive_body"].(schema.DynamicAttribute)
	if !ok {
		t.Fatalf("expected sensitive_body to be a dynamic attribute")
	}
	if !sensitiveBody.WriteOnly {
		t.Fatalf("expected sensitive_body to be write-only")
	}

	sensitiveBodyVersion, ok := response.Schema.Attributes["sensitive_body_version"].(schema.MapAttribute)
	if !ok {
		t.Fatalf("expected sensitive_body_version to be a map attribute")
	}
	if sensitiveBodyVersion.ElementType != types.StringType {
		t.Fatalf("expected sensitive_body_version element type to be string")
	}
	if response.Schema.Version != 3 {
		t.Fatalf("expected schema version 3, got %d", response.Schema.Version)
	}
}

func TestActionEphemeralSchemaSensitiveBody(t *testing.T) {
	var response fwephemeral.SchemaResponse
	r := ActionEphemeral{}
	r.Schema(context.Background(), fwephemeral.SchemaRequest{}, &response)

	sensitiveBody, ok := response.Schema.Attributes["sensitive_body"].(ephemeralschema.DynamicAttribute)
	if !ok {
		t.Fatalf("expected sensitive_body to be a dynamic attribute")
	}
	if sensitiveBody.IsWriteOnly() {
		t.Fatalf("expected ephemeral sensitive_body not to be write-only")
	}
}

func TestBuildActionRequestBodyMergesSensitiveBody(t *testing.T) {
	body := mustDynamicFromJSON(t, `{
		"source": {
			"registryUri": "dhi.io",
			"sourceImage": "node:22.22.1-alpine3.23"
		},
		"targetTags": ["base/dhi-node"]
	}`)
	sensitiveBody := mustDynamicFromJSON(t, `{
		"source": {
			"credentials": {
				"username": "user",
				"password": "secret"
			}
		}
	}`)

	actual, err := buildActionRequestBody(body, sensitiveBody, types.MapNull(types.StringType), types.MapNull(types.StringType))
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	assertJsonEqual(t, actual, `{
		"source": {
			"registryUri": "dhi.io",
			"sourceImage": "node:22.22.1-alpine3.23",
			"credentials": {
				"username": "user",
				"password": "secret"
			}
		},
		"targetTags": ["base/dhi-node"]
	}`)
}

func TestBuildActionRequestBodySkipsUnchangedSensitiveBodyVersion(t *testing.T) {
	body := mustDynamicFromJSON(t, `{
		"properties": {
			"mode": "NoForce"
		}
	}`)
	sensitiveBody := mustDynamicFromJSON(t, `{
		"properties": {
			"sharedKey": "secret"
		}
	}`)
	version := types.MapValueMust(types.StringType, map[string]attr.Value{
		"properties.sharedKey": types.StringValue("v1"),
	})

	actual, err := buildActionRequestBody(body, sensitiveBody, version, version)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	assertJsonEqual(t, actual, `{
		"properties": {
			"mode": "NoForce"
		}
	}`)
}

func mustDynamicFromJSON(t *testing.T, input string) types.Dynamic {
	t.Helper()
	value, err := dynamic.FromJSONImplied([]byte(input))
	if err != nil {
		t.Fatalf("failed to build dynamic value: %+v", err)
	}
	return value
}

func assertJsonEqual(t *testing.T, actual interface{}, expected string) {
	t.Helper()

	var expectedValue interface{}
	if err := json.Unmarshal([]byte(expected), &expectedValue); err != nil {
		t.Fatalf("failed to unmarshal expected JSON: %+v", err)
	}

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("failed to marshal actual value: %+v", err)
	}
	var actualValue interface{}
	if err := json.Unmarshal(actualJSON, &actualValue); err != nil {
		t.Fatalf("failed to unmarshal actual JSON: %+v", err)
	}

	if !reflect.DeepEqual(actualValue, expectedValue) {
		t.Fatalf("expected %s but got %s", mustMarshalJSON(expectedValue), mustMarshalJSON(actualValue))
	}
}

func mustMarshalJSON(value interface{}) string {
	out, _ := json.Marshal(value)
	return string(out)
}
