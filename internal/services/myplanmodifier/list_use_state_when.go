package myplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ListSemanticallyEqualFunc func(a, b types.List) bool

func ListUseStateWhen(equalFunc ListSemanticallyEqualFunc) planmodifier.List {
	return listUseStateWhen{
		EqualFunc: equalFunc,
	}
}

type listUseStateWhen struct {
	EqualFunc ListSemanticallyEqualFunc
}

func (u listUseStateWhen) Description(ctx context.Context) string {
	return "Use the state value when new value is functionally equivalent to the old and thus no change is required."
}

func (u listUseStateWhen) MarkdownDescription(ctx context.Context) string {
	return "Use the state value when new value is functionally equivalent to the old and thus no change is required."
}

func (u listUseStateWhen) PlanModifyList(ctx context.Context, request planmodifier.ListRequest, response *planmodifier.ListResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	for _, element := range request.ConfigValue.Elements() {
		if element.IsUnknown() {
			return
		}
	}
	if request.StateValue.IsNull() || request.StateValue.IsUnknown() {
		return
	}
	if u.EqualFunc(request.ConfigValue, request.StateValue) {
		response.PlanValue = request.StateValue
	}
}
