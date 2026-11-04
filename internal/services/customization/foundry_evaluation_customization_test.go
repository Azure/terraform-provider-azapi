package customization

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

func TestFoundryEvaluationCustomizationCreateResult(t *testing.T) {
	transport := &foundryTestTransport{
		responses: []foundryTestResponse{
			{
				statusCode: http.StatusCreated,
				body:       `{"id":"eval_123","name":"evaluation"}`,
			},
		},
	}
	client := newFoundryTestClient(t, transport)
	inputID := parse.DataPlaneResourceId{
		AzureResourceType: "Microsoft.Foundry/evaluation/versions",
		ApiVersion:        "2025-05-01",
		Name:              "__generated__",
		ParentId:          "account.services.ai.azure.com/api/projects/project",
	}

	createdID, responseBody, err := (FoundryEvaluationCustomization{}).CreateResultFunc()(
		t.Context(),
		client,
		inputID,
		map[string]interface{}{
			"name":               "evaluation",
			"data_source_config": map[string]interface{}{"type": "custom"},
		},
		clients.RequestOptions{
			Headers:         map[string]string{"X-Test": "header"},
			QueryParameters: map[string]string{"custom": "value"},
		},
	)
	if err != nil {
		t.Fatalf("creating evaluation: %v", err)
	}

	if createdID.AzureResourceId != "account.services.ai.azure.com/api/projects/project/openai/v1/evals/eval_123" {
		t.Fatalf("unexpected evaluation resource ID: %q", createdID.AzureResourceId)
	}
	if createdID.Name != "eval_123" {
		t.Fatalf("unexpected evaluation name: %q", createdID.Name)
	}
	if responseBody == nil {
		t.Fatal("expected response body")
	}

	request := transport.lastRequest(t)
	if request.method != http.MethodPost {
		t.Fatalf("unexpected method: %q", request.method)
	}
	if request.url != "https://account.services.ai.azure.com/api/projects/project/openai/v1/evals?custom=value" {
		t.Fatalf("unexpected URL: %q", request.url)
	}
	if request.headers.Get("X-Test") != "header" {
		t.Fatalf("expected custom header, got %q", request.headers.Get("X-Test"))
	}
	if strings.Contains(request.url, "api-version") {
		t.Fatalf("unexpected API version in URL: %q", request.url)
	}

	var requestBody map[string]interface{}
	if err := json.Unmarshal([]byte(request.body), &requestBody); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if requestBody["name"] != "evaluation" {
		t.Fatalf("unexpected request body: %#v", requestBody)
	}
}

func TestFoundryEvaluationCustomizationUpdateBody(t *testing.T) {
	body, err := foundryEvaluationUpdateRequestBody(map[string]interface{}{
		"name":               "updated",
		"metadata":           map[string]interface{}{"team": "platform"},
		"data_source_config": map[string]interface{}{"type": "custom"},
		"testing_criteria":   []interface{}{"immutable"},
	})
	if err != nil {
		t.Fatalf("building update body: %v", err)
	}

	expected := map[string]interface{}{
		"name":     "updated",
		"metadata": map[string]interface{}{"team": "platform"},
	}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("unexpected update body: got %#v, want %#v", body, expected)
	}
}

func TestFoundryEvaluationCustomizationResponseValidation(t *testing.T) {
	tests := []struct {
		name     string
		response interface{}
	}{
		{name: "nil", response: nil},
		{name: "missing ID", response: map[string]interface{}{}},
		{name: "non-string ID", response: map[string]interface{}{"id": 123}},
		{name: "empty ID", response: map[string]interface{}{"id": " "}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := foundryResourceID(test.response); err == nil {
				t.Fatal("expected response validation error")
			}
		})
	}
}

func TestFoundryEvaluationCustomizationLifecycle(t *testing.T) {
	customization := FoundryEvaluationCustomization{}

	if customization.GetResourceType() != "Microsoft.Foundry/evaluation/versions" {
		t.Fatalf("unexpected resource type: %q", customization.GetResourceType())
	}
	if customization.CreateFunc() != nil {
		t.Fatal("evaluation customization should use CreateResultFunc")
	}
	if customization.CreateResultFunc() == nil ||
		customization.ReadFunc() == nil ||
		customization.UpdateFunc() == nil ||
		customization.DeleteFunc() == nil {
		t.Fatal("evaluation customization must define complete lifecycle behavior")
	}
}

type foundryTestResponse struct {
	statusCode int
	body       string
}

type foundryTestRequest struct {
	method  string
	url     string
	headers http.Header
	body    string
}

type foundryTestTransport struct {
	responses []foundryTestResponse
	requests  []foundryTestRequest
}

func (t *foundryTestTransport) Do(request *http.Request) (*http.Response, error) {
	if len(t.responses) == 0 {
		return nil, fmt.Errorf("unexpected %s request", request.Method)
	}

	var body []byte
	if request.Body != nil {
		var err error
		body, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
	}
	t.requests = append(t.requests, foundryTestRequest{
		method:  request.Method,
		url:     request.URL.String(),
		headers: request.Header.Clone(),
		body:    string(body),
	})

	response := t.responses[0]
	t.responses = t.responses[1:]
	return &http.Response{
		StatusCode: response.statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(response.body)),
		Request:    request,
	}, nil
}

func (transport *foundryTestTransport) lastRequest(t testing.TB) foundryTestRequest {
	t.Helper()
	if len(transport.requests) == 0 {
		t.Fatal("expected a request")
	}
	return transport.requests[len(transport.requests)-1]
}

type foundryTestCredential struct{}

func (foundryTestCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     "token",
		ExpiresOn: time.Now().Add(time.Hour),
	}, nil
}

func newFoundryTestClient(t *testing.T, transport policy.Transporter) clients.Client {
	t.Helper()

	dataPlaneClient, err := clients.NewDataPlaneClient(
		foundryTestCredential{},
		&arm.ClientOptions{
			ClientOptions: policy.ClientOptions{Transport: transport},
		},
	)
	if err != nil {
		t.Fatalf("creating data-plane client: %v", err)
	}

	return clients.Client{DataPlaneClient: dataPlaneClient}
}
