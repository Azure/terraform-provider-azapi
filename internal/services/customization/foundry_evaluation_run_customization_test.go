package customization

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

func TestFoundryEvaluationRunCustomizationCreateResult(t *testing.T) {
	transport := &foundryTestTransport{
		responses: []foundryTestResponse{
			{
				statusCode: http.StatusCreated,
				body:       `{"id":"run_456","status":"queued"}`,
			},
		},
	}
	client := newFoundryTestClient(t, transport)
	inputID := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Foundry/evaluation/runs",
		ApiVersion:        "2025-05-01",
		Name:              "__generated__",
		ParentId:          "account.services.ai.azure.com/api/projects/project",
		EvaluationId:      "eval_123",
	}

	createdID, _, err := (FoundryEvaluationRunCustomization{}).CreateResultFunc()(
		t.Context(),
		client,
		inputID,
		map[string]interface{}{
			"name":          "run",
			"evaluation_id": "must-not-be-sent",
			"data_source":   map[string]interface{}{"type": "azure_ai_target_completions"},
		},
		clients.RequestOptions{},
	)
	if err != nil {
		t.Fatalf("creating evaluation run: %v", err)
	}

	if createdID.AzureResourceId != "account.services.ai.azure.com/api/projects/project/openai/v1/evals/eval_123/runs/run_456" {
		t.Fatalf("unexpected evaluation run resource ID: %q", createdID.AzureResourceId)
	}
	if createdID.Name != "run_456" || createdID.EvaluationId != "eval_123" {
		t.Fatalf("unexpected generated ID fields: %#v", createdID)
	}

	request := transport.lastRequest(t)
	if request.method != http.MethodPost {
		t.Fatalf("unexpected method: %q", request.method)
	}
	if request.url != "https://account.services.ai.azure.com/api/projects/project/openai/v1/evals/eval_123/runs" {
		t.Fatalf("unexpected URL: %q", request.url)
	}
	if strings.Contains(request.url, "api-version") {
		t.Fatalf("unexpected API version in URL: %q", request.url)
	}

	var requestBody map[string]interface{}
	if err := json.Unmarshal([]byte(request.body), &requestBody); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if _, exists := requestBody["evaluation_id"]; exists {
		t.Fatalf("evaluation_id must not be sent in request body: %#v", requestBody)
	}
	if requestBody["name"] != "run" {
		t.Fatalf("unexpected request body: %#v", requestBody)
	}
}

func TestFoundryEvaluationRunCustomizationValidation(t *testing.T) {
	t.Run("requires evaluation ID", func(t *testing.T) {
		_, err := foundryEvaluationRunEvaluationID(parse.DataPlaneResourceId{})
		if err == nil || !strings.Contains(err.Error(), "evaluation_id") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("accepts trimmed evaluation ID", func(t *testing.T) {
		id, err := foundryEvaluationRunEvaluationID(parse.DataPlaneResourceId{EvaluationId: " eval_123 "})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != "eval_123" {
			t.Fatalf("unexpected evaluation ID: %q", id)
		}
	})

	t.Run("removes both evaluation ID body spellings", func(t *testing.T) {
		body, err := foundryEvaluationRunRequestBody(map[string]interface{}{
			"evaluation_id": "one",
			"evaluationId":  "two",
			"name":          "run",
		})
		if err != nil {
			t.Fatalf("building request body: %v", err)
		}
		if _, exists := body["evaluation_id"]; exists {
			t.Fatal("snake-case evaluation_id was not removed")
		}
		if _, exists := body["evaluationId"]; exists {
			t.Fatal("camel-case evaluationId was not removed")
		}
		if body["name"] != "run" {
			t.Fatalf("unexpected request body: %#v", body)
		}
	})
}

func TestFoundryEvaluationRunIDFromResponse(t *testing.T) {
	tests := []struct {
		name     string
		response interface{}
		expected string
	}{
		{name: "id", response: map[string]interface{}{"id": "run_123"}, expected: "run_123"},
		{name: "run_id", response: map[string]interface{}{"run_id": "run_123"}, expected: "run_123"},
		{name: "runId", response: map[string]interface{}{"runId": "run_123"}, expected: "run_123"},
		{name: "runID", response: map[string]interface{}{"runID": "run_123"}, expected: "run_123"},
		{name: "text", response: "run_123", expected: "run_123"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := foundryEvaluationRunIDFromResponse(test.response)
			if err != nil {
				t.Fatalf("extracting run ID: %v", err)
			}
			if actual != test.expected {
				t.Fatalf("unexpected run ID: %q", actual)
			}
		})
	}

	for _, response := range []interface{}{
		nil,
		map[string]interface{}{},
		map[string]interface{}{"id": 123},
		map[string]interface{}{"id": " "},
	} {
		t.Run("invalid response", func(t *testing.T) {
			if _, err := foundryEvaluationRunIDFromResponse(response); err == nil {
				t.Fatal("expected response validation error")
			}
		})
	}
}

func TestFoundryEvaluationRunCustomizationLifecycle(t *testing.T) {
	customization := FoundryEvaluationRunCustomization{}

	if customization.GetResourceType() != "Microsoft.Foundry/evaluation/runs" {
		t.Fatalf("unexpected resource type: %q", customization.GetResourceType())
	}
	if customization.CreateFunc() != nil {
		t.Fatal("evaluation run customization should use CreateResultFunc")
	}
	if customization.CreateResultFunc() == nil ||
		customization.ReadFunc() == nil ||
		customization.UpdateFunc() == nil ||
		customization.DeleteFunc() == nil {
		t.Fatal("evaluation run customization must define complete lifecycle behavior")
	}

	err := customization.UpdateFunc()(
		t.Context(),
		clients.Client{},
		parse.DataPlaneResourceId{},
		nil,
		clients.RequestOptions{},
	)
	if err == nil || !strings.Contains(err.Error(), "cannot be updated") {
		t.Fatalf("unexpected update error: %v", err)
	}
}
