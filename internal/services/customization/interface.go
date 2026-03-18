package customization

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
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

type ReadFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error)
type DeleteFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error
type CreateFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error
type CreateResultFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) (parse.DataPlaneResourceId, interface{}, error)
type UpdateFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error
