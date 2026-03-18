package planmodifierdynamic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplaceIfFunc is a conditional function used in the RequiresReplaceIf
// plan modifier to determine whether the attribute requires replacement.
type RequiresReplaceIfFunc func(context.Context, planmodifier.DynamicRequest, *RequiresReplaceIfFuncResponse)

func RequiresReplace() planmodifier.Dynamic {
	return RequiresReplaceIf(
		func(_ context.Context, _ planmodifier.DynamicRequest, resp *RequiresReplaceIfFuncResponse) {
			resp.RequiresReplace = true
		},
		"If the value of this attribute changes, Terraform will destroy and recreate the resource.",
		"If the value of this attribute changes, Terraform will destroy and recreate the resource.",
	)
}

// RequiresReplaceIfNotNull is a plan modifier that sets RequiresReplace
// on the attribute if the planned value is different from the state value.
// It will not trigger replacement when either the planned or state value is effectively null.
// A Dynamic attribute can be null in two ways:
//  1. The outer Dynamic wrapper is null (e.g., types.DynamicNull())
//  2. The outer Dynamic is "known" but wraps a null underlying value
//     (e.g., types.DynamicValue(TupleNull())) — this happens when HCL evaluates
//     a conditional like `condition ? [...] : null` for a DynamicAttribute.
//
// When this function is called, the equality of the planned and state values
// has already been checked by the PlanModifyDynamic method, so we can assume they are different values.
func RequiresReplaceIfNotNull() planmodifier.Dynamic {
	return RequiresReplaceIf(
		func(_ context.Context, req planmodifier.DynamicRequest, resp *RequiresReplaceIfFuncResponse) {
			planNull := req.PlanValue.IsNull() || req.PlanValue.UnderlyingValue().IsNull()
			stateNull := req.StateValue.IsNull() || req.StateValue.UnderlyingValue().IsNull()
			resp.RequiresReplace = !planNull && !stateNull
		},
		"If the planned value is different from the state value, Terraform will destroy and recreate the resource unless either the planned or state value is null.",
		"If the planned value is different from the state value, Terraform will destroy and recreate the resource unless either the planned or state value is null.",
	)
}

// RequiresReplaceIfFuncResponse is the response type for a RequiresReplaceIfFunc.
type RequiresReplaceIfFuncResponse struct {
	// Diagnostics report errors or warnings related to this logic. An empty
	// or unset slice indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics

	// RequiresReplace should be enabled if the resource should be replaced.
	RequiresReplace bool
}

func RequiresReplaceIf(f RequiresReplaceIfFunc, description, markdownDescription string) planmodifier.Dynamic {
	return requiresReplaceIfModifier{
		ifFunc:              f,
		description:         description,
		markdownDescription: markdownDescription,
	}
}

// requiresReplaceIfModifier is an plan modifier that sets RequiresReplace
// on the attribute if a given function is true.
type requiresReplaceIfModifier struct {
	ifFunc              RequiresReplaceIfFunc
	description         string
	markdownDescription string
}

// Description returns a human-readable description of the plan modifier.
func (m requiresReplaceIfModifier) Description(_ context.Context) string {
	return m.description
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m requiresReplaceIfModifier) MarkdownDescription(_ context.Context) string {
	return m.markdownDescription
}

// PlanModifyString implements the plan modification logic.
func (m requiresReplaceIfModifier) PlanModifyDynamic(ctx context.Context, req planmodifier.DynamicRequest, resp *planmodifier.DynamicResponse) {
	// Do not replace on resource creation.
	if req.State.Raw.IsNull() {
		return
	}

	// Do not replace on resource destroy.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Do not replace if the plan and state values are equal.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	ifFuncResp := &RequiresReplaceIfFuncResponse{}

	m.ifFunc(ctx, req, ifFuncResp)

	resp.Diagnostics.Append(ifFuncResp.Diagnostics...)
	resp.RequiresReplace = ifFuncResp.RequiresReplace
}
