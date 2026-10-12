package customization

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataPlaneResource interface {
	GetResourceType() string
	CreateFunc() CreateFunc
	// CreateResultFunc is optional. When implemented, it enables create flows where the service
	// generates the final resource identifier (for example, POST to a collection endpoint that
	// returns an "id"). The returned DataPlaneResourceId will be used for subsequent read/state.
	//
	// If CreateResultFunc is non-nil, azapi_data_plane_resource will prefer it over CreateFunc.
	//
	// Note: This is intentionally an optional interface (via type assertion) to avoid forcing
	// all existing customizations to implement it.
	ReadFunc() ReadFunc
	UpdateFunc() UpdateFunc
	DeleteFunc() DeleteFunc
}

// DataPlaneResourceWithCreateResult is an optional extension interface for customizations that
// need to return the server-generated ID from create.
type DataPlaneResourceWithCreateResult interface {
	DataPlaneResource
	CreateResultFunc() CreateResultFunc
}

// DataPlaneResourceWithCreateResponse is an optional extension for customizations
// that need to pass the create response into post-create output augmentation.
// Unlike DataPlaneResourceWithCreateResult, it does not imply a generated ID.
type DataPlaneResourceWithCreateResponse interface {
	DataPlaneResource
	CreateResponseFunc() CreateResponseFunc
}

// DataPlaneResourceWithPlanBody is an optional extension for customizations
// which need to carry provider-computed body values from state into the plan.
// Its Terraform values preserve unknowns while the plan body is normalized.
type DataPlaneResourceWithPlanBody interface {
	DataPlaneResource
	PlanBodyFunc() PlanBodyFunc
}

// DataPlaneResourceWithReadOptions is an optional extension for customizations
// whose read response must not be merged into the request body.
type DataPlaneResourceWithReadOptions interface {
	DataPlaneResource
	PreserveBodyStateOnRead() bool
	UseResponseBodyAsOutput() bool
	AugmentReadOutput(
		ctx context.Context,
		responseBody interface{},
		stateBody interface{},
		createResponse interface{},
	) (interface{}, error)
}

type ReadFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error)
type DeleteFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error
type CreateFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error
type CreateResultFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) (parse.DataPlaneResourceId, interface{}, error)
type CreateResponseFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) (interface{}, error)
type PlanBodyFunc = func(ctx context.Context, planBody types.Dynamic, stateBody types.Dynamic) (types.Dynamic, error)
type UpdateFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error
