package clients

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type dataPlaneTestCredential struct{}

func (dataPlaneTestCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "test-token"}, nil
}

type dataPlaneTestTransport struct {
	request *http.Request
}

func (t *dataPlaneTestTransport) Do(request *http.Request) (*http.Response, error) {
	t.request = request
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		Request:    request,
	}, nil
}

func TestActionWithContentType(t *testing.T) {
	transport := &dataPlaneTestTransport{}
	client, err := NewDataPlaneClient(dataPlaneTestCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud:     cloud.AzurePublic,
			Transport: transport,
		},
	})
	if err != nil {
		t.Fatalf("building data plane client: %v", err)
	}

	_, err = client.ActionWithContentType(
		context.Background(),
		"example.com/resource",
		"",
		"2025-05-01",
		http.MethodPut,
		map[string]string{"name": "example"},
		RequestOptions{},
		"application/vnd.azure+json",
	)
	if err != nil {
		t.Fatalf("calling action: %v", err)
	}

	if transport.request == nil {
		t.Fatal("expected action request")
	}
	if got := transport.request.Header.Get("Content-Type"); got != "application/vnd.azure+json" {
		t.Fatalf("expected content type application/vnd.azure+json, got %q", got)
	}
	if got := transport.request.URL.Query().Get("api-version"); got != "2025-05-01" {
		t.Fatalf("expected API version 2025-05-01, got %q", got)
	}
	if got := transport.request.URL.Path; got != "/resource" {
		t.Fatalf("expected resource path /resource, got %q", got)
	}
}
