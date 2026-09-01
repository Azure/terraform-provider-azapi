package services

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_readOverrideFromObject(t *testing.T) {
	value := types.ObjectValueMust(readOverrideAttributeTypes(), map[string]attr.Value{
		"method":        types.StringValue("POST"),
		"action":        types.StringValue("list"),
		"response_path": types.StringValue("value"),
	})

	model, diags := readOverrideFromObject(context.Background(), value)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if model.Method.ValueString() != "POST" || model.Action.ValueString() != "list" || model.ResponsePath.ValueString() != "value" {
		t.Fatalf("unexpected read override: %#v", model)
	}

	model, diags = readOverrideFromObject(context.Background(), types.ObjectNull(readOverrideAttributeTypes()))
	if diags.HasError() || model != nil {
		t.Fatalf("null read override should decode to nil without diagnostics")
	}
}

func Test_resolveReadResponseWithoutOverride(t *testing.T) {
	var getResponse interface{}
	_ = json.Unmarshal([]byte(`{"properties":{},"name":"appsettings"}`), &getResponse)

	result, diags := resolveReadResponse(context.Background(), nil, "", "", getResponse, nil, clients.RequestOptions{})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !reflect.DeepEqual(result, getResponse) {
		t.Fatalf("expected the GET response to be returned unchanged")
	}
}

func Test_applyResponsePath(t *testing.T) {
	var response interface{}
	_ = json.Unmarshal([]byte(`{"value":{"properties":{"MY_SETTING":"value1"}}}`), &response)

	var diags diag.Diagnostics
	result := applyResponsePath(response, "value", &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	var expected interface{}
	_ = json.Unmarshal([]byte(`{"properties":{"MY_SETTING":"value1"}}`), &expected)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %v but got %v", expected, result)
	}

	result = applyResponsePath(response, "", &diags)
	if !reflect.DeepEqual(result, response) {
		t.Fatalf("empty path should return the response unchanged")
	}
}
