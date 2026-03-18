package customization

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

// FoundryAgentCustomization supports Foundry agents which use
// POST-to-collection create and named-resource get/update/delete.
//
// Create: POST {endpoint}/agents?api-version=v1
// Get:    GET  {endpoint}/agents/{agentName}?api-version=v1
// Update: POST {endpoint}/agents/{agentName}?api-version=v1
// Delete: DELETE {endpoint}/agents/{agentName}?api-version=v1
type FoundryAgentCustomization struct{}

func (c FoundryAgentCustomization) GetResourceType() string {
	return "Microsoft.Foundry/agents"
}

func (c FoundryAgentCustomization) CreateFunc() CreateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		collectionId := strings.TrimSuffix(id.ParentId, "/") + "/agents"
		_, err := client.DataPlaneClient.Action(ctx, collectionId, "", id.ApiVersion, http.MethodPost, body, options)
		return err
	}
}

func (c FoundryAgentCustomization) ReadFunc() ReadFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
		return client.DataPlaneClient.Get(ctx, id, options)
	}
}

func (c FoundryAgentCustomization) UpdateFunc() UpdateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "", id.ApiVersion, http.MethodPost, body, options)
		return err
	}
}

func (c FoundryAgentCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "", id.ApiVersion, http.MethodDelete, nil, options)
		return err
	}
}

var _ DataPlaneResource = &FoundryAgentCustomization{}
