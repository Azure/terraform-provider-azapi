package clients

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"testing"
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
				"User-Agent": defaultUserAgent,
			},
		},
		{
			Name: "Request with SuppliedUserAgentHeader",
			RequestHeaders: map[string]string{
				"Supplied-User-Agent": "suppliedUserAgent",
			},
			ExpectedHeaders: map[string]string{
				"User-Agent": defaultUserAgent + " suppliedUserAgent",
			},
		},
		{
			Name: "Request with UserAgent",
			RequestHeaders: map[string]string{
				"User-Agent": "customUserAgent",
			},
			ExpectedHeaders: map[string]string{
				"User-Agent": "customUserAgent",
			},
		},
		{
			Name: "Request with UserAgent and SuppliedUserAgentHeader",
			RequestHeaders: map[string]string{
				"User-Agent":          "customUserAgent",
				"Supplied-User-Agent": "suppliedUserAgent",
			},
			ExpectedHeaders: map[string]string{
				"User-Agent": "customUserAgent suppliedUserAgent",
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
