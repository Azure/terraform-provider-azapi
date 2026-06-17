package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	HeaderPolicyExternalEvaluations    = "x-ms-policy-external-evaluations"
	acquirePolicyTokenAPIVersion       = "2025-03-01"
	errorCodeRequestDisallowedByPolicy = "RequestDisallowedByPolicy"
	additionalInfoTypePolicyViolation  = "PolicyViolation"
)

type acquirePolicyTokenPolicy struct {
	pipeline       runtime.Pipeline
	endpoint       string
	subscriptionID string
}

var _ policy.Policy = &acquirePolicyTokenPolicy{}

type policyErrorResponse struct {
	Error policyError `json:"error"`
}

type policyError struct {
	Code           string                      `json:"code"`
	AdditionalInfo []policyErrorAdditionalInfo `json:"additionalInfo"`
}

type policyErrorAdditionalInfo struct {
	Type string              `json:"type"`
	Info policyViolationInfo `json:"info"`
}

type policyViolationInfo struct {
	EvaluationDetails policyEvaluationDetails `json:"evaluationDetails"`
}

type policyEvaluationDetails struct {
	MissingPolicyTokenDetails *missingPolicyTokenDetails `json:"missingPolicyTokenDetails"`
}

type missingPolicyTokenDetails struct {
	ShouldDeny                bool   `json:"shouldDeny"`
	IsChangeReferenceRequired bool   `json:"isChangeReferenceRequired"`
	EndpointKind              string `json:"endpointKind"`
}

func NewAcquirePolicyTokenPolicy(pipeline runtime.Pipeline, endpoint string, subscriptionID string) policy.Policy {
	return &acquirePolicyTokenPolicy{
		pipeline:       pipeline,
		endpoint:       endpoint,
		subscriptionID: subscriptionID,
	}
}

func (p *acquirePolicyTokenPolicy) Do(req *policy.Request) (*http.Response, error) {
	if req.Raw().Method == http.MethodGet {
		return req.Next()
	}

	if req.Raw().Header.Get(HeaderPolicyExternalEvaluations) != "" {
		return req.Next()
	}

	resp, err := req.Next()
	if err != nil {
		return resp, err
	}

	if !runtime.HasStatusCode(resp, http.StatusForbidden) {
		return resp, nil
	}

	details, err := missingPolicyTokenDetailsFromResponse(resp)
	if err != nil || details == nil {
		return resp, nil
	}

	if details.IsChangeReferenceRequired {
		return nil, fmt.Errorf("the request was disallowed by an invoke policy that requires a change reference, but change reference is not yet supported by this provider")
	}

	token, err := p.acquirePolicyToken(req)
	if err != nil {
		return nil, fmt.Errorf("the request was disallowed by an invoke policy and acquiring a policy token to satisfy it failed: %w (original policy error: %w)", err, runtime.NewResponseError(resp))
	}
	if token == "" {
		return resp, nil
	}

	if err := req.RewindBody(); err != nil {
		return nil, fmt.Errorf("rewinding request body: %w", err)
	}
	req.Raw().Header.Set(HeaderPolicyExternalEvaluations, token)
	return req.Next()
}

func missingPolicyTokenDetailsFromResponse(resp *http.Response) (*missingPolicyTokenDetails, error) {
	body, err := runtime.Payload(resp)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var errResp policyErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, nil
	}
	if !strings.EqualFold(errResp.Error.Code, errorCodeRequestDisallowedByPolicy) {
		return nil, nil
	}

	for i := range errResp.Error.AdditionalInfo {
		info := errResp.Error.AdditionalInfo[i]
		if !strings.EqualFold(info.Type, additionalInfoTypePolicyViolation) {
			continue
		}
		details := info.Info.EvaluationDetails.MissingPolicyTokenDetails
		if details == nil {
			continue
		}
		if !details.ShouldDeny {
			continue
		}
		return details, nil
	}
	return nil, nil
}

func (p *acquirePolicyTokenPolicy) acquirePolicyToken(req *policy.Request) (string, error) {
	rawReq := req.Raw()

	var content interface{}
	if body := req.Body(); body != nil {
		if _, err := body.Seek(0, io.SeekStart); err != nil {
			return "", fmt.Errorf("seeking request body before reading to acquire policy token: %w", err)
		}
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			return "", fmt.Errorf("reading request body: %w", err)
		}
		if _, err := body.Seek(0, io.SeekStart); err != nil {
			return "", fmt.Errorf("seeking request body after reading to acquire policy token: %w", err)
		}
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &content); err != nil {
				content = string(bodyBytes)
			}
		}
	}

	requestBody := map[string]interface{}{
		"operation": map[string]interface{}{
			"uri":        rawReq.URL.String(),
			"httpMethod": rawReq.Method,
			"content":    content,
		},
	}

	ctx := rawReq.Context()
	urlPath := fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/acquirePolicyToken", p.subscriptionID)
	acquireReq, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(p.endpoint, urlPath))
	if err != nil {
		return "", fmt.Errorf("creating acquire policy token request: %w", err)
	}

	reqQP := acquireReq.Raw().URL.Query()
	reqQP.Set("api-version", acquirePolicyTokenAPIVersion)
	acquireReq.Raw().URL.RawQuery = reqQP.Encode()
	acquireReq.Raw().Header.Set("Accept", "application/json")
	acquireReq.Raw().Header.Set("x-ms-force-sync", "true")
	if err := runtime.MarshalAsJSON(acquireReq, requestBody); err != nil {
		return "", fmt.Errorf("marshalling acquire policy token request body: %w", err)
	}

	resp, err := p.pipeline.Do(acquireReq)
	if err != nil {
		return "", err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return "", runtime.NewResponseError(resp)
	}

	var responseBody struct {
		Result string `json:"result"`
		Token  string `json:"token"`
	}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return "", fmt.Errorf("unmarshalling acquire policy token response: %w", err)
	}
	if !strings.EqualFold(responseBody.Result, "Succeeded") || responseBody.Token == "" {
		return "", nil
	}
	return responseBody.Token, nil
}
