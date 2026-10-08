package clients

import (
	"context"
	"net/http"
	"testing"
)

func TestBuildRequestAPIParameters(t *testing.T) {
	t.Run("original builder preserves empty API version", func(t *testing.T) {
		request, err := buildRequest(
			context.Background(),
			RequestOptions{
				QueryParameters: map[string]string{"foo": "bar"},
			},
			"https://example.com/resource?existing=value",
			http.MethodGet,
			"",
		)
		if err != nil {
			t.Fatalf("building request: %v", err)
		}

		query := request.Raw().URL.Query()
		if _, exists := query["api-version"]; !exists {
			t.Fatal("expected original builder to include the API version parameter")
		}
		if got := query.Get("existing"); got != "value" {
			t.Fatalf("expected existing query parameter, got %q", got)
		}
		if got := query.Get("foo"); got != "bar" {
			t.Fatalf("expected custom query parameter, got %q", got)
		}
	})

	t.Run("Foundry builder omits API version", func(t *testing.T) {
		request, err := buildRequestWithoutAPIVersion(
			context.Background(),
			RequestOptions{},
			"https://example.com/resource",
			http.MethodGet,
			"",
		)
		if err != nil {
			t.Fatalf("building request: %v", err)
		}

		if _, exists := request.Raw().URL.Query()["api-version"]; exists {
			t.Fatal("expected Foundry builder to omit the API version parameter")
		}
	})

	t.Run("sets API version", func(t *testing.T) {
		request, err := buildRequest(
			context.Background(),
			RequestOptions{},
			"https://example.com/resource",
			http.MethodGet,
			"2025-05-01",
		)
		if err != nil {
			t.Fatalf("building request: %v", err)
		}

		if got := request.Raw().URL.Query().Get("api-version"); got != "2025-05-01" {
			t.Fatalf("expected API version, got %q", got)
		}
	})
}
