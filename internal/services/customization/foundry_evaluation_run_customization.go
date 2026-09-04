package customization

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

// FoundryEvaluationRunCustomization manages Foundry evaluation runs.
//
// The top-level evaluation_id identifies the evaluation collection and is
// available through id.EvaluationId. It is used in the request URL and is not
// sent in the request body.
//
// Foundry generates the run identifier, so the top-level name argument must
// not be configured.
type FoundryEvaluationRunCustomization struct{}

const foundryEvaluationRunResourceType = "Microsoft.Foundry/evaluation/runs"

func (c FoundryEvaluationRunCustomization) GetResourceType() string {
	return foundryEvaluationRunResourceType
}

func foundryEvaluationRunCollectionID(
	parentID string,
	evaluationID string,
) string {
	return fmt.Sprintf(
		"%s/openai/v1/evals/%s/runs",
		strings.TrimRight(strings.TrimSpace(parentID), "/"),
		url.PathEscape(strings.TrimSpace(evaluationID)),
	)
}

func foundryEvaluationBodyMap(
	body interface{},
) (map[string]interface{}, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf(
			"marshalling evaluation body: %w",
			err,
		)
	}

	var values map[string]interface{}
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf(
			"unmarshalling evaluation body: %w",
			err,
		)
	}

	if values == nil {
		return nil, fmt.Errorf(
			"evaluation body must be an object",
		)
	}

	return values, nil
}

func foundryEvaluationRunRequestBody(
	body interface{},
) (map[string]interface{}, error) {
	values, err := foundryEvaluationBodyMap(body)
	if err != nil {
		return nil, err
	}

	// Remove defensively in case the generic resource layer injects this field
	// into the body. The configured top-level evaluation_id is obtained from
	// id.EvaluationId.
	delete(values, "evaluation_id")
	delete(values, "evaluationId")

	return values, nil
}

func foundryEvaluationRunIDFromResponse(
	response interface{},
) (string, error) {
	if value, ok := response.(string); ok {
		runID := strings.TrimSpace(value)
		if runID != "" {
			return runID, nil
		}
	}

	values, err := foundryEvaluationBodyMap(response)
	if err != nil {
		return "", err
	}

	for _, fieldName := range []string{
		"id",
		"run_id",
		"runId",
		"runID",
	} {
		value, exists := values[fieldName]
		if !exists || value == nil {
			continue
		}

		runID, ok := value.(string)
		if !ok {
			return "", fmt.Errorf(
				"evaluation run response field %q must be a string",
				fieldName,
			)
		}

		runID = strings.TrimSpace(runID)
		if runID == "" {
			return "", fmt.Errorf(
				"evaluation run response field %q must not be empty",
				fieldName,
			)
		}

		return runID, nil
	}

	return "", fmt.Errorf(
		"evaluation run response did not contain a generated run ID",
	)
}

func foundryEvaluationRunEvaluationID(
	id parse.DataPlaneResourceId,
) (string, error) {
	evaluationID := strings.TrimSpace(id.EvaluationId)

	if evaluationID == "" || evaluationID == "__generated__" {
		return "", fmt.Errorf(
			`top-level "evaluation_id" is required for %s`,
			foundryEvaluationRunResourceType,
		)
	}

	return evaluationID, nil
}

func (c FoundryEvaluationRunCustomization) CreateFunc() CreateFunc {
	// Foundry generates the run ID. Creation uses CreateResultFunc.
	return nil
}

func (c FoundryEvaluationRunCustomization) CreateResultFunc() CreateResultFunc {
	return func(
		ctx context.Context,
		client clients.Client,
		id parse.DataPlaneResourceId,
		body interface{},
		options clients.RequestOptions,
	) (parse.DataPlaneResourceId, interface{}, error) {
		evaluationID, err := foundryEvaluationRunEvaluationID(id)
		if err != nil {
			return parse.DataPlaneResourceId{}, nil, err
		}

		requestBody, err := foundryEvaluationRunRequestBody(body)
		if err != nil {
			return parse.DataPlaneResourceId{}, nil, err
		}

		responseBody, err := client.DataPlaneClient.ActionWithoutAPIVersion(
			ctx,
			foundryEvaluationRunCollectionID(id.ParentId, evaluationID),
			"",
			http.MethodPost,
			requestBody,
			options,
		)
		if err != nil {
			return parse.DataPlaneResourceId{}, nil, err
		}

		runID, err := foundryEvaluationRunIDFromResponse(responseBody)
		if err != nil {
			return parse.DataPlaneResourceId{}, responseBody, err
		}

		createdID, err := parse.NewDataPlaneResourceIdWithEvaluationID(
			runID,
			id.ParentId,
			evaluationID,
			fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion),
		)
		if err != nil {
			return parse.DataPlaneResourceId{}, responseBody, err
		}

		return createdID, responseBody, nil
	}
}

func (c FoundryEvaluationRunCustomization) ReadFunc() ReadFunc {
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

func (c FoundryEvaluationRunCustomization) UpdateFunc() UpdateFunc {
	return func(
		_ context.Context,
		_ clients.Client,
		_ parse.DataPlaneResourceId,
		_ interface{},
		_ clients.RequestOptions,
	) error {
		return fmt.Errorf(
			"Foundry evaluation runs cannot be updated; create a new run instead",
		)
	}
}

func (c FoundryEvaluationRunCustomization) DeleteFunc() DeleteFunc {
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

var _ DataPlaneResource = &FoundryEvaluationRunCustomization{}
var _ DataPlaneResourceWithCreateResult = &FoundryEvaluationRunCustomization{}
