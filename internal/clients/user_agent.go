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
	userAgent := c.UserAgent
	// if the request already has a User-Agent header, append it to the provider's default User-Agent
	if requestUserAgent := req.Raw().Header.Get(HeaderUserAgent); requestUserAgent != "" {
		userAgent = userAgent + " " + requestUserAgent
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
