package clients

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type rawTestCredential struct {
	mu     sync.Mutex
	scopes []string
}

func (c *rawTestCredential) GetToken(_ context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.mu.Lock()
	c.scopes = append(c.scopes, options.Scopes...)
	c.mu.Unlock()
	return azcore.AccessToken{Token: "token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type rawTestTransport struct {
	mu       sync.Mutex
	requests []*http.Request
	bodies   [][]byte
	calls    int
}

func (t *rawTestTransport) Do(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	t.mu.Lock()
	t.requests = append(t.requests, req.Clone(req.Context()))
	t.bodies = append(t.bodies, body)
	t.calls++
	call := t.calls
	t.mu.Unlock()

	status := http.StatusCreated
	if call == 1 {
		status = http.StatusServiceUnavailable
	}
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    req,
	}, nil
}

func TestDataPlaneClientDoRawReplaysBodyAndUsesVersionHeader(t *testing.T) {
	credential := &rawTestCredential{}
	transport := &rawTestTransport{}
	config := cloud.Configuration{
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			"Storage": {
				Endpoint: "https://.blob.core.windows.net",
				Audience: "https://storage.azure.com",
			},
		},
	}
	client, err := NewDataPlaneClient(credential, &arm.ClientOptions{ClientOptions: policy.ClientOptions{
		Cloud:     config,
		Transport: transport,
		Retry: policy.RetryOptions{
			MaxRetries: 1,
			RetryDelay: time.Millisecond,
		},
	}})
	if err != nil {
		t.Fatal(err)
	}

	content := []byte{0, 1, 2, '"', '\n'}
	_, err = client.DoRaw(context.Background(), http.MethodPut, "account.blob.core.windows.net/container/a b/#c?d/%", RawRequestOptions{
		RequestOptions: RequestOptions{
			Headers: map[string]string{"Content-Type": "application/octet-stream"},
		},
		APIVersionHeader: "2023-11-03",
		Content:          content,
		SuccessCodes:     []int{http.StatusCreated},
	})
	if err != nil {
		t.Fatal(err)
	}

	if transport.calls != 2 {
		t.Fatalf("expected 2 requests, got %d", transport.calls)
	}
	for i, body := range transport.bodies {
		if string(body) != string(content) {
			t.Fatalf("request %d body differs: got %v, want %v", i, body, content)
		}
	}
	request := transport.requests[1]
	if got := request.URL.EscapedPath(); got != "/container/a%20b/%23c%3Fd/%25" {
		t.Fatalf("unexpected escaped path %q", got)
	}
	if request.URL.Query().Has("api-version") {
		t.Fatalf("raw request unexpectedly has api-version query parameter")
	}
	if got := request.Header.Get("x-ms-version"); got != "2023-11-03" {
		t.Fatalf("unexpected x-ms-version %q", got)
	}
	if got := request.Header.Get("Content-Type"); got != "application/octet-stream" {
		t.Fatalf("unexpected Content-Type %q", got)
	}
	if request.ContentLength != int64(len(content)) {
		t.Fatalf("unexpected Content-Length %d", request.ContentLength)
	}
	if len(credential.scopes) != 1 || credential.scopes[0] != "https://storage.azure.com/.default" {
		t.Fatalf("unexpected token scopes %v", credential.scopes)
	}
}

func TestDataPlaneClientDoRawRejectsUnconfiguredService(t *testing.T) {
	client, err := NewDataPlaneClient(&rawTestCredential{}, &arm.ClientOptions{ClientOptions: policy.ClientOptions{
		Cloud: cloud.Configuration{
			Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
				cloud.ResourceManager: {
					Endpoint: "https://management.example.com",
					Audience: "https://management.example.com",
				},
			},
		},
		Transport: &rawTestTransport{},
	}})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.DoRaw(context.Background(), http.MethodHead, "account.blob.custom.example/container/blob", RawRequestOptions{
		APIVersionHeader: "2023-11-03",
	})
	if err == nil || !strings.Contains(err.Error(), "unsupported data-plane endpoint") {
		t.Fatalf("expected unsupported endpoint error, got %v", err)
	}
}

type rawSuccessTransport struct{}

func (rawSuccessTransport) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    req,
	}, nil
}

func TestDataPlaneClientDoRawSelectsStorageScope(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		endpoint string
		audience string
	}{
		{"public", "account.blob.core.windows.net", "https://.blob.core.windows.net", "https://storage.azure.com"},
		{"government", "account.blob.core.usgovcloudapi.net", "https://.blob.core.usgovcloudapi.net", "https://storage.azure.us"},
		{"china", "account.blob.core.chinacloudapi.cn", "https://.blob.core.chinacloudapi.cn", "https://storage.azure.cn"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			credential := &rawTestCredential{}
			client, err := NewDataPlaneClient(credential, &arm.ClientOptions{ClientOptions: policy.ClientOptions{
				Cloud: cloud.Configuration{
					Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
						"Storage": {Endpoint: test.endpoint, Audience: test.audience},
					},
				},
				Transport: rawSuccessTransport{},
			}})
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.DoRaw(context.Background(), http.MethodHead, test.host+"/container/blob", RawRequestOptions{
				APIVersionHeader: "2023-11-03",
			})
			if err != nil {
				t.Fatal(err)
			}
			_ = resp.Body.Close()
			if len(credential.scopes) != 1 || credential.scopes[0] != test.audience+"/.default" {
				t.Fatalf("unexpected token scopes %v", credential.scopes)
			}
		})
	}
}

func TestBuildEscapedURLPreservesBlobName(t *testing.T) {
	tests := map[string]string{
		"account.blob.core.windows.net/container/a b":         "https://account.blob.core.windows.net/container/a%20b",
		"account.blob.core.windows.net/container/a#b?c%d":     "https://account.blob.core.windows.net/container/a%23b%3Fc%25d",
		"account.blob.core.windows.net/container/目录/文件":       "https://account.blob.core.windows.net/container/%E7%9B%AE%E5%BD%95/%E6%96%87%E4%BB%B6",
		"account.blob.core.windows.net/container/a//b/.././c": "https://account.blob.core.windows.net/container/a//b/.././c",
	}
	for resourceID, expected := range tests {
		actual, err := buildEscapedURL(resourceID)
		if err != nil {
			t.Fatalf("building URL for %q: %v", resourceID, err)
		}
		if actual != expected {
			t.Fatalf("unexpected URL for %q: got %q, want %q", resourceID, actual, expected)
		}
	}
}
