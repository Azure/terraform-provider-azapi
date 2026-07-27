package clients

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// sigintTransport simulates ARM while an azapi_resource Create is interrupted by
// SIGINT, reproducing https://github.com/Azure/terraform-provider-azapi/issues/1110.
//
// Sequence modelled:
//  1. PUT (CreateOrUpdate) is accepted by ARM: the resource is created
//     server-side and a long-running-operation (LRO) poller is returned via the
//     Azure-AsyncOperation header.
//  2. While the provider polls the LRO, terraform core receives SIGINT and
//     cancels the context (simulated here by cancelling on the first poll GET).
//     The poll therefore fails with context.Canceled.
//
// Like the real net/http transport, Do returns ctx.Err() when the request's
// context is already cancelled. This is what makes the provider's post-failure
// recovery Get in azapi_resource.go (which reuses the same cancelled context)
// fail as well, leaving the resource orphaned in Azure but absent from state.
type sigintTransport struct {
	host         string
	resourcePath string
	asyncOpPath  string
	cancel       context.CancelFunc

	mu        sync.Mutex
	created   bool
	pollCount int
}

func (t *sigintTransport) Do(req *http.Request) (*http.Response, error) {
	if err := req.Context().Err(); err != nil {
		return nil, err
	}

	switch {
	case req.Method == http.MethodPut && req.URL.Path == t.resourcePath:
		// ARM accepts the create and starts an async operation.
		t.mu.Lock()
		t.created = true
		t.mu.Unlock()
		resp := &http.Response{
			StatusCode: http.StatusCreated,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(
				`{"id":"` + t.resourcePath + `","properties":{"provisioningState":"Creating"}}`)),
			Request: req,
		}
		resp.Header.Set("Azure-AsyncOperation", t.host+t.asyncOpPath)
		return resp, nil

	case req.Method == http.MethodGet && req.URL.Path == t.asyncOpPath:
		// SIGINT arrives mid-poll: cancel the context exactly as terraform core
		// would when the operation is stopped, then report "still in progress".
		// azcore observes the cancellation and PollUntilDone returns context.Canceled.
		t.mu.Lock()
		t.pollCount++
		t.mu.Unlock()
		if t.cancel != nil {
			t.cancel()
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"status":"InProgress"}`)),
			Request:    req,
		}, nil

	case req.Method == http.MethodGet && req.URL.Path == t.resourcePath:
		// A read of the resource: ARM already created it, so it is returned with
		// its id and a terminal provisioning state.
		t.mu.Lock()
		created := t.created
		t.mu.Unlock()
		if !created {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"error":{"code":"ResourceNotFound"}}`)),
				Request:    req,
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(
				`{"id":"` + t.resourcePath + `","properties":{"provisioningState":"Succeeded"}}`)),
			Request: req,
		}, nil
	}

	return nil, fmt.Errorf("unexpected request %s %s", req.Method, req.URL.Path)
}

func (t *sigintTransport) wasCreated() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.created
}

// TestResourceClientCreateOrUpdate_SIGINTLeavesResourceOrphaned reproduces
// https://github.com/Azure/terraform-provider-azapi/issues/1110 and validates the
// approach used by the fix.
//
// When an in-flight Create is interrupted by SIGINT, ARM has already created the
// resource but the provider cancels its poll context and CreateOrUpdate returns
// context.Canceled. A recovery Get that reuses the SAME cancelled context also
// fails, so nothing is written to state and the resource is orphaned in Azure.
//
// The fix in azapi_resource.go CreateUpdate performs the recovery Get on a
// context detached from cancellation (context.WithoutCancel + a bounded timeout),
// which — as asserted below — successfully reads back the id ARM created so it can
// be persisted to state.
func TestResourceClientCreateOrUpdate_SIGINTLeavesResourceOrphaned(t *testing.T) {
	const (
		host         = "https://management.azure.com"
		resourcePath = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Storage/storageAccounts/stcanceltestsample"
		asyncOpPath  = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/locations/westeurope/asyncoperations/11111111-1111-1111-1111-111111111111"
		apiVersion   = "2023-01-01"
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transport := &sigintTransport{
		host:         host,
		resourcePath: resourcePath,
		asyncOpPath:  asyncOpPath,
		cancel:       cancel,
	}

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{}, &policy.ClientOptions{
		Telemetry: policy.TelemetryOptions{Disabled: true},
		Transport: transport,
	})
	client := &ResourceClient{
		host: host,
		pl:   pl,
	}

	// 1. In-flight Create interrupted by SIGINT. ARM accepted the PUT (the
	//    resource now exists) but the poll context is cancelled, so the provider's
	//    CreateOrUpdate fails with context.Canceled and returns no id.
	body := map[string]interface{}{
		"location": "westeurope",
		"sku":      map[string]interface{}{"name": "Standard_LRS"},
		"kind":     "StorageV2",
	}
	if _, err := client.CreateOrUpdate(ctx, resourcePath, apiVersion, body, RequestOptions{}); err == nil {
		t.Fatal("expected CreateOrUpdate to fail once the poll context was cancelled")
	} else if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled from the interrupted poll, got %v", err)
	}

	if !transport.wasCreated() {
		t.Fatal("precondition failed: ARM should have accepted the PUT and created the resource")
	}
	if transport.pollCount == 0 {
		t.Fatal("precondition failed: the LRO poller should have been reached before cancellation")
	}

	// 2. The naive recovery (the pre-fix behaviour) reuses the SAME context that
	//    SIGINT cancelled, so the recovery Get fails, no id is captured, and
	//    terraform writes nothing to state even though the resource exists in Azure.
	if _, err := client.Get(ctx, resourcePath, apiVersion, RequestOptions{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected recovery Get on the cancelled context to fail with context.Canceled, got %v", err)
	}

	// 3. The fix: azapi_resource.go CreateUpdate runs the recovery Get on a context
	//    detached from cancellation, which reads back the resource ARM created and
	//    yields the id the provider persists to state, avoiding the orphan.
	recoverCtx := context.WithoutCancel(ctx)
	responseBody, err := client.Get(recoverCtx, resourcePath, apiVersion, RequestOptions{})
	if err != nil {
		t.Fatalf("recovery Get on a detached context should succeed, got %v", err)
	}
	responseMap, ok := responseBody.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map response body, got %T", responseBody)
	}
	if responseMap["id"] != resourcePath {
		t.Fatalf("expected recovered resource id %q, got %v", resourcePath, responseMap["id"])
	}
}
