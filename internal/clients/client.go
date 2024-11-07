package clients

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/features"
)

type Client struct {
	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Features features.UserFeatures

	ResourceClient  *ResourceClient
	DataPlaneClient *DataPlaneClient

	Account ResourceManagerAccount
}

type Option struct {
	Cred                        azcore.TokenCredential
	ApplicationUserAgent        string
	Features                    features.UserFeatures
	SkipProviderRegistration    bool
	DisableCorrelationRequestID bool
	CloudCfg                    cloud.Configuration
	CustomCorrelationRequestID  string
	SubscriptionId              string
	TenantId                    string
}

// NOTE: it should be possible for this method to become Private once the top level Client's removed

func (client *Client) Build(ctx context.Context, o *Option) error {
	client.StopContext = ctx
	client.Features = o.Features

	azlog.SetListener(func(cls azlog.Event, msg string) {
		log.Printf("[DEBUG] %s %s: %s\n", time.Now().Format(time.StampMicro), cls, msg)
	})

	perCallPolicies := make([]policy.Policy, 0)
	perCallPolicies = append(perCallPolicies, withUserAgent(o.ApplicationUserAgent))
	if !o.DisableCorrelationRequestID {
		id := o.CustomCorrelationRequestID
		if id == "" {
			id = correlationRequestID()
		}
		perCallPolicies = append(perCallPolicies, withCorrelationRequestID(id))
	}
	perRetryPolicies := make([]policy.Policy, 0)
	perRetryPolicies = append(perRetryPolicies, NewLiveTrafficLogPolicy())

	allowedHeaders := []string{
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"Elapsed-Time",
		"Location",
		"Metadata",
		"Ocp-Automation-Accountid",
		"P3p",
		"Strict-Transport-Security",
		"Vary",
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-Ms-Correlation-Request-Id",
		"X-Ms-Ests-Server",
		"X-Ms-Failure-Cause",
		"X-Ms-Ratelimit-Remaining-Subscription-Reads",
		"X-Ms-Ratelimit-Remaining-Subscription-Writes",
		"X-Ms-Ratelimit-Remaining-Tenant-Reads",
		"X-Ms-Ratelimit-Remaining-Tenant-Writes",
		"X-Ms-Request-Id",
		"X-Ms-Routing-Request-Id",
		"X-Xss-Protection",
	}
	allowedQueryParams := []string{
		"api-version",
		"$skipToken",
	}
	tok, err := o.Cred.GetToken(client.StopContext, policy.TokenRequestOptions{
		TenantID: o.TenantId,
		Scopes:   []string{o.CloudCfg.Services[cloud.ResourceManager].Endpoint + "/.default"}})
	if err != nil {
		return fmt.Errorf("failed to get token to determine object id: %w", err)
	}
	cl, err := parseTokenClaims(tok.Token)
	if err != nil {
		return fmt.Errorf("failed to parse token claims to determine object id: %w", err)
	}
	resourceClient, err := NewResourceClient(o.Cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: o.CloudCfg,
			// Disable the default telemetry policy, because it has a length limitation for user agent
			Telemetry: policy.TelemetryOptions{
				Disabled: true,
			},
			Logging: policy.LogOptions{
				IncludeBody:        true,
				AllowedHeaders:     allowedHeaders,
				AllowedQueryParams: allowedQueryParams,
			},
			PerCallPolicies:  perCallPolicies,
			PerRetryPolicies: perRetryPolicies,
		},
		DisableRPRegistration: o.SkipProviderRegistration,
	})
	if err != nil {
		return err
	}
	client.ResourceClient = resourceClient

	dataPlaneClient, err := NewDataPlaneClient(o.Cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: o.CloudCfg,
			// Disable the default telemetry policy, because it has a length limitation for user agent
			Telemetry: policy.TelemetryOptions{
				Disabled: true,
			},
			Logging: policy.LogOptions{
				IncludeBody:        true,
				AllowedHeaders:     allowedHeaders,
				AllowedQueryParams: allowedQueryParams,
			},
			PerCallPolicies:  perCallPolicies,
			PerRetryPolicies: perRetryPolicies,
		},
		DisableRPRegistration: o.SkipProviderRegistration,
	})
	if err != nil {
		return err
	}
	client.DataPlaneClient = dataPlaneClient

	var oid string
	if cl != nil && cl.ObjectId != "" {
		oid = cl.ObjectId
	}
	client.Account = NewResourceManagerAccount(o.TenantId, o.SubscriptionId, oid)

	return nil
}

func parseTokenClaims(token string) (*tokenClaims, error) {
	// Parse the token to get the claims
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("parseTokenClaims: token does not have 3 parts")
	}
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("parseTokenClaims: error decoding token: %s", err)
	}
	var claims tokenClaims
	err = json.Unmarshal(decoded, &claims)
	if err != nil {
		return nil, fmt.Errorf("parseTokenClaims: error unmarshalling claims: %w", err)
	}
	return &claims, nil
}

type tokenClaims struct {
	Audience          string   `json:"aud"`
	Expires           int64    `json:"exp"`
	IssuedAt          int64    `json:"iat"`
	Issuer            string   `json:"iss"`
	IdentityProvider  string   `json:"idp"`
	ObjectId          string   `json:"oid"`
	Roles             []string `json:"roles"`
	Scopes            string   `json:"scp"`
	Subject           string   `json:"sub"`
	TenantRegionScope string   `json:"tenant_region_scope"`
	TenantId          string   `json:"tid"`
	Version           string   `json:"ver"`

	AppDisplayName string `json:"app_displayname,omitempty"`
	AppId          string `json:"appid,omitempty"`
	IdType         string `json:"idtyp,omitempty"`
}
