package clients

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func TestApplyAPIVersion_DefaultsToQueryParameter(t *testing.T) {
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://example.table.core.windows.net/Tables('test')")
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	applyAPIVersion(req, "2026-04-06", RequestOptions{})

	if got := req.Raw().URL.Query().Get("api-version"); got != "2026-04-06" {
		t.Fatalf("expected api-version query parameter, got %q", got)
	}
	if got := req.Raw().Header.Get("x-ms-version"); got != "" {
		t.Fatalf("expected no x-ms-version header, got %q", got)
	}
}

func TestApplyAPIVersion_CanUseHeaderInsteadOfQueryParameter(t *testing.T) {
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, "https://example.table.core.windows.net/Tables('test')")
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	applyAPIVersion(req, "2026-04-06", RequestOptions{
		DisableAPIVersionQueryParameter: true,
		APIVersionHeaderName:            "x-ms-version",
	})

	if got := req.Raw().URL.Query().Get("api-version"); got != "" {
		t.Fatalf("expected no api-version query parameter, got %q", got)
	}
	if got := req.Raw().Header.Get("x-ms-version"); got != "2026-04-06" {
		t.Fatalf("expected x-ms-version header, got %q", got)
	}
}
