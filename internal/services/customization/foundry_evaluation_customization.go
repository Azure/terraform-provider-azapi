package customization

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

// FoundryEvaluationCustomization supports Foundry evaluations, which use a
// collection POST for creation and a generated evaluation ID for later calls.
//
// Create: POST {endpoint}/openai/v1/evals
// Get:    GET  {endpoint}/openai/v1/evals/{evalID}
// Update: POST {endpoint}/openai/v1/evals/{evalID}
// Delete: DELETE {endpoint}/openai/v1/evals/{evalID}
type FoundryEvaluationCustomization struct{}

const foundryEvaluationCollectionPath = "/openai/v1/evals"

func (c FoundryEvaluationCustomization) GetResourceType() string {
	return "Microsoft.Foundry/evaluation/versions"
}

func (c FoundryEvaluationCustomization) CreateFunc() CreateFunc {
	return nil
}

func foundryEvaluationCollectionID(parentID string) string {
	return strings.TrimRight(strings.TrimSpace(parentID), "/") + foundryEvaluationCollectionPath
}

func foundryResourceID(response interface{}) (string, error) {
	values, err := foundryEvaluationBodyMap(response)
	if err != nil {
		return "", fmt.Errorf("invalid evaluation response: %w", err)
	}

	value, ok := values["id"]
	if !ok {
		return "", fmt.Errorf(`evaluation response field "id" is missing`)
	}
	evaluationID, ok := value.(string)
	if !ok {
		return "", fmt.Errorf(`evaluation response field "id" must be a string`)
	}
	if strings.TrimSpace(evaluationID) == "" {
		return "", fmt.Errorf(`evaluation response field "id" must be a non-empty string`)
	}

	return strings.TrimSpace(evaluationID), nil
}

func foundryEvaluationUpdateRequestBody(body interface{}) (map[string]interface{}, error) {
	values, err := foundryEvaluationBodyMap(body)
	if err != nil {
		return nil, err
	}

	// The update API only accepts mutable evaluation metadata. The data source
	// and testing criteria are create-time properties.
	updateBody := make(map[string]interface{}, 2)
	for _, fieldName := range []string{"name", "metadata"} {
		if value, exists := values[fieldName]; exists {
			updateBody[fieldName] = value
		}
	}

	return updateBody, nil
}

func (c FoundryEvaluationCustomization) CreateResultFunc() CreateResultFunc {
	return func(
		ctx context.Context,
		client clients.Client,
		id parse.DataPlaneResourceId,
		body interface{},
		options clients.RequestOptions,
	) (parse.DataPlaneResourceId, interface{}, error) {
		responseBody, err := client.DataPlaneClient.ActionWithoutAPIVersion(
			ctx,
			foundryEvaluationCollectionID(id.ParentId),
			"",
			http.MethodPost,
			body,
			options,
		)
		if err != nil {
			return parse.DataPlaneResourceId{}, nil, err
		}

		evaluationID, err := foundryResourceID(responseBody)
		if err != nil {
			return parse.DataPlaneResourceId{}, nil, err
		}

		createdID, err := parse.NewDataPlaneResourceId(
			evaluationID,
			id.ParentId,
			fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion),
		)
		if err != nil {
			return parse.DataPlaneResourceId{}, nil, err
		}

		return createdID, responseBody, nil
	}
}

func (c FoundryEvaluationCustomization) ReadFunc() ReadFunc {
	return func(
		ctx context.Context,
		client clients.Client,
		id parse.DataPlaneResourceId,
		options clients.RequestOptions,
	) (interface{}, error) {
		return client.DataPlaneClient.ActionWithoutAPIVersion(
			ctx,
			id.AzureResourceId,
			"",
			http.MethodGet,
			nil,
			options,
		)
	}
}

func (c FoundryEvaluationCustomization) UpdateFunc() UpdateFunc {
	return func(
		ctx context.Context,
		client clients.Client,
		id parse.DataPlaneResourceId,
		body interface{},
		options clients.RequestOptions,
	) error {
		requestBody, err := foundryEvaluationUpdateRequestBody(body)
		if err != nil {
			return err
		}

		_, err = client.DataPlaneClient.ActionWithoutAPIVersion(
			ctx,
			id.AzureResourceId,
			"",
			http.MethodPost,
			requestBody,
			options,
		)
		return err
	}
}

func (c FoundryEvaluationCustomization) DeleteFunc() DeleteFunc {
	return func(
		ctx context.Context,
		client clients.Client,
		id parse.DataPlaneResourceId,
		options clients.RequestOptions,
	) error {
		_, err := client.DataPlaneClient.ActionWithoutAPIVersion(
			ctx,
			id.AzureResourceId,
			"",
			http.MethodDelete,
			nil,
			options,
		)
		return err
	}
}

var _ DataPlaneResource = &FoundryEvaluationCustomization{}
var _ DataPlaneResourceWithCreateResult = &FoundryEvaluationCustomization{}
