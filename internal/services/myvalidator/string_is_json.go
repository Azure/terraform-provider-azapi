package myvalidator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringIsJSON struct{}

func (v stringIsJSON) Description(ctx context.Context) string {
	return "validate this in JSON format"
}

func (v stringIsJSON) MarkdownDescription(ctx context.Context) string {
	return "validate this in JSON format"
}

func (_ stringIsJSON) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	str := req.ConfigValue

	if str.IsUnknown() || str.IsNull() {
		return
	}

	if _, errs := StringIsJSONPluginSDK(str.ValueString(), req.Path.String()); len(errs) != 0 {
		for _, err := range errs {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid JSON string",
				err.Error())
		}
	}
}

func StringIsJSON() validator.String {
	return stringIsJSON{}
}

// StringIsJSONPluginSDK is a SchemaValidateFunc which tests to make sure the supplied string is valid JSON.
func StringIsJSONPluginSDK(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	if _, err := NormalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}

	return warnings, errors
}

// NormalizeJsonString Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func NormalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	bytes, _ := json.Marshal(j)
	return string(bytes), nil
}
