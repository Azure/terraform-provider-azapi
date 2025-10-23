package customization

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

type DataPlaneResource interface {
	GetResourceType() string
	CreateFunc() CreateFunc
	ReadFunc() ReadFunc
	UpdateFunc() UpdateFunc
	DeleteFunc() DeleteFunc
}

type ReadFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error)
type DeleteFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error
type CreateFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error
type UpdateFunc = func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error
