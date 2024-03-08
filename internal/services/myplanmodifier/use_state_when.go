package myplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StringSemanticallyEqualFunc func(a, b types.String) bool

func UseStateWhen(equalFunc StringSemanticallyEqualFunc) planmodifier.String {
	return useStateWhen{
		EqualFunc: equalFunc,
	}
}

type useStateWhen struct {
	EqualFunc StringSemanticallyEqualFunc
}

func (u useStateWhen) Description(ctx context.Context) string {
	return "Use the state value when new value is functionally equivalent to the old and thus no change is required."
}

func (u useStateWhen) MarkdownDescription(ctx context.Context) string {
	return "Use the state value when new value is functionally equivalent to the old and thus no change is required."
}

func (u useStateWhen) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
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
