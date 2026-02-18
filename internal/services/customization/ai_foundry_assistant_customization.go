package customization

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

// AIFoundryAssistantCustomization supports AI Foundry Agents "assistants" which use
// POST-to-collection create and return a server-generated identifier.
//
// Create: POST {endpoint}/assistants?api-version=v1
// Get:    GET  {endpoint}/assistants/{assistantId}?api-version=v1
// Update: POST {endpoint}/assistants/{assistantId}?api-version=v1
// Delete: DELETE {endpoint}/assistants/{assistantId}?api-version=v1
type AIFoundryAssistantCustomization struct{}

func (c AIFoundryAssistantCustomization) GetResourceType() string {
	return "Microsoft.AIFoundry/agents/assistants"
}

func (c AIFoundryAssistantCustomization) CreateFunc() CreateFunc {
	// Prefer CreateResultFunc for this customization.
	return nil
}

func (c AIFoundryAssistantCustomization) CreateResultFunc() CreateResultFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) (parse.DataPlaneResourceId, interface{}, error) {
		collectionId := strings.TrimSuffix(id.ParentId, "/") + "/assistants"
		resp, err := client.DataPlaneClient.Action(ctx, collectionId, "", id.ApiVersion, http.MethodPost, body, options)
		if err != nil {
			return parse.DataPlaneResourceId{}, nil, err
		}

		assistantId, err := extractStringField(resp, "id")
		if err != nil {
			return parse.DataPlaneResourceId{}, resp, err
		}

		resourceType := fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion)
		newId, err := parse.NewDataPlaneResourceId(assistantId, id.ParentId, resourceType)
		if err != nil {
			return parse.DataPlaneResourceId{}, resp, err
		}

		return newId, resp, nil
	}
}

func (c AIFoundryAssistantCustomization) ReadFunc() ReadFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
		return client.DataPlaneClient.Get(ctx, id, options)
	}
}

func (c AIFoundryAssistantCustomization) UpdateFunc() UpdateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "", id.ApiVersion, http.MethodPost, body, options)
		return err
	}
}

func (c AIFoundryAssistantCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "", id.ApiVersion, http.MethodDelete, nil, options)
		return err
	}
}

var _ DataPlaneResource = &AIFoundryAssistantCustomization{}
var _ DataPlaneResourceWithCreateResult = &AIFoundryAssistantCustomization{}

func extractStringField(payload interface{}, field string) (string, error) {
	if payload == nil {
		return "", fmt.Errorf("response body is nil")
	}
	obj, ok := payload.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("response body is not an object")
	}
	value, ok := obj[field]
	if !ok {
		return "", fmt.Errorf("response body missing field %q", field)
	}
	stringValue, ok := value.(string)
	if !ok || strings.TrimSpace(stringValue) == "" {
		return "", fmt.Errorf("response body field %q is not a non-empty string", field)
	}
	return stringValue, nil
}
