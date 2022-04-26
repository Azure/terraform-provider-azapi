package clients

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	HeaderUserAgent = "User-Agent"
)

type UserAgentPolicy struct {
	UserAgent string
}

func (c UserAgentPolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Set(HeaderUserAgent, c.UserAgent)
	return req.Next()
}

var _ policy.Policy = UserAgentPolicy{}

// withUserAgent returns a policy.Policy that adds an HTTP extension header of
// `User-Agent` whose value is passed and has no length limitation
func withUserAgent(userAgent string) policy.Policy {
	return UserAgentPolicy{UserAgent: userAgent}
}
