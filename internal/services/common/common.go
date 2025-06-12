package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func AsStringList(input types.List) []string {
	var result []string
	diags := input.ElementsAs(context.Background(), &result, false)
	if diags.HasError() {
		tflog.Warn(context.Background(), fmt.Sprintf("failed to convert list to string list: %s", diags))
	}
	return result
}

func AsMapOfString(input types.Map) map[string]string {
	result := make(map[string]string)
	diags := input.ElementsAs(context.Background(), &result, false)
	if diags.HasError() {
		tflog.Warn(context.Background(), fmt.Sprintf("failed to convert input to map of strings: %s", diags))
	}
	return result
}

func AsMapOfLists(input types.Map) map[string][]string {
	result := make(map[string][]string)
	diags := input.ElementsAs(context.Background(), &result, false)
	if diags.HasError() {
		tflog.Warn(context.Background(), fmt.Sprintf("failed to convert input to map of lists: %s", diags))
	}
	return result
}
