package myplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MapSemanticallyEqualFunc func(a, b types.Map) bool

func MapUseStateWhen(equalFunc MapSemanticallyEqualFunc) planmodifier.Map {
	return mapUseStateWhen{
		EqualFunc: equalFunc,
	}
}

type mapUseStateWhen struct {
	EqualFunc MapSemanticallyEqualFunc
}

func (u mapUseStateWhen) Description(ctx context.Context) string {
	return "Use the state value when new value is functionally equivalent to the old and thus no change is required."
}

func (u mapUseStateWhen) MarkdownDescription(ctx context.Context) string {
	return "Use the state value when new value is functionally equivalent to the old and thus no change is required."
}

func (u mapUseStateWhen) PlanModifyMap(ctx context.Context, request planmodifier.MapRequest, response *planmodifier.MapResponse) {
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
