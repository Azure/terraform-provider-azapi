package services

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStorageTableEntityResponseStateShaping(t *testing.T) {
	responseBody := map[string]interface{}{
		"PartitionKey": "pk",
		"RowKey":       "rk",
		"Timestamp":    "2026-09-13T00:00:00Z",
		"odata.etag":   `W/"datetime'2026-09-13T00%3A00%3A00.0000000Z'"`,
		"outputs":      "remote-value",
	}

	t.Run("server-only identity and metadata stay out of body", func(t *testing.T) {
		configuredBody := map[string]interface{}{
			"outputs": "configured-value",
		}
		expectedBody := map[string]interface{}{
			"outputs": "remote-value",
		}
		managedBody := utils.UpdateObject(configuredBody, responseBody, utils.UpdateJsonOption{
			IgnoreMissingProperty: true,
		})
		if !reflect.DeepEqual(managedBody, expectedBody) {
			t.Fatalf("expected managed body %#v, got %#v", expectedBody, managedBody)
		}
	})

	exportAll := types.DynamicValue(types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("*"),
	}))
	output, err := buildOutputFromBody(responseBody, exportAll, nil)
	if err != nil {
		t.Fatalf("building output: %v", err)
	}
	outputValue, ok := output.UnderlyingValue().(types.Dynamic)
	if !ok {
		t.Fatalf("expected dynamic output value, got %T", output.UnderlyingValue())
	}
	outputJSON, err := dynamic.ToJSON(outputValue)
	if err != nil {
		t.Fatalf("converting output to JSON: %v", err)
	}
	var actualOutput map[string]interface{}
	if err := json.Unmarshal(outputJSON, &actualOutput); err != nil {
		t.Fatalf("decoding output: %v", err)
	}
	if !reflect.DeepEqual(actualOutput, responseBody) {
		t.Fatalf("expected full response output %#v, got %#v", responseBody, actualOutput)
	}
}
