package clients

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	HeaderUserAgent         = "User-Agent"
	HeaderSuppliedUserAgent = "Supplied-User-Agent"
)

type UserAgentPolicy struct {
	UserAgent string
}

func (c UserAgentPolicy) Do(req *policy.Request) (*http.Response, error) {
	userAgent := req.Raw().Header.Get(HeaderUserAgent)
	if userAgent == "" {
		userAgent = c.UserAgent
	}

	if suppliedUserAgent := req.Raw().Header.Get(HeaderSuppliedUserAgent); suppliedUserAgent != "" {
		userAgent = userAgent + " " + suppliedUserAgent
		req.Raw().Header.Del(HeaderSuppliedUserAgent)
	}

	req.Raw().Header.Set(HeaderUserAgent, userAgent)
	return req.Next()
}

var _ policy.Policy = UserAgentPolicy{}

// withUserAgent returns a policy.Policy that adds an HTTP extension header of
// `User-Agent` whose value is passed and has no length limitation
func withUserAgent(userAgent string) policy.Policy {
	return UserAgentPolicy{UserAgent: userAgent}
}
