package clients

import (
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func Test_UserAgentPolicy(t *testing.T) {
	const defaultUserAgent = "defaultUserAgent"
	testcases := []struct {
		Name            string
		RequestHeaders  map[string]string
		ExpectedHeaders map[string]string
	}{
		{
			Name:           "Request without UserAgent",
			RequestHeaders: map[string]string{},
			ExpectedHeaders: map[string]string{
				HeaderUserAgent: defaultUserAgent,
			},
		},
		{
			Name: "Request with UserAgent",
			RequestHeaders: map[string]string{
				HeaderUserAgent: "customUserAgent",
			},
			ExpectedHeaders: map[string]string{
				HeaderUserAgent: defaultUserAgent + " customUserAgent",
			},
		},
	}

	userAgentPolicy := withUserAgent(defaultUserAgent)

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			rawRequest := &http.Request{
				Header: make(http.Header),
			}
			for k, v := range tc.RequestHeaders {
				rawRequest.Header.Set(k, v)
			}

			req, _ := runtime.NewRequestFromRequest(rawRequest)
			_, _ = userAgentPolicy.Do(req)

			for k, v := range tc.ExpectedHeaders {
				if actualV := req.Raw().Header.Get(k); actualV != v {
					t.Errorf("Expected header %s to be %s, got %s", k, v, actualV)
				}
			}
		})
	}

}
