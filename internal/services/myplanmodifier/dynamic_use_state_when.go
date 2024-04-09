package myplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DynamicSemanticallyEqualFunc func(a, b types.Dynamic) bool

func DynamicUseStateWhen(equalFunc DynamicSemanticallyEqualFunc) planmodifier.Dynamic {
	return dynamicUseStateWhen{
		EqualFunc: equalFunc,
	}
}

type dynamicUseStateWhen struct {
	EqualFunc DynamicSemanticallyEqualFunc
}

func (u dynamicUseStateWhen) Description(ctx context.Context) string {
	return "Use the state value when new value is functionally equivalent to the old and thus no change is required."
}

func (u dynamicUseStateWhen) MarkdownDescription(ctx context.Context) string {
	return "Use the state value when new value is functionally equivalent to the old and thus no change is required."
}

func (u dynamicUseStateWhen) PlanModifyDynamic(ctx context.Context, request planmodifier.DynamicRequest, response *planmodifier.DynamicResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if request.StateValue.IsNull() || request.StateValue.IsUnknown() {
		return
	}
	if u.EqualFunc(request.ConfigValue, request.StateValue) {
		response.PlanValue = request.StateValue
	}
}
