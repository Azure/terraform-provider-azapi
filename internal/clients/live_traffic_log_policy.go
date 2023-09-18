package clients

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const redactedValue = "REDACTED"

type liveTrafficLogPolicy struct {
	notAllowedHeaders map[string]bool
}

type traffic struct {
	LiveRequest  liveRequest  `json:"request"`
	LiveResponse liveResponse `json:"response"`
}

type liveRequest struct {
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Body    string            `json:"body"`
}

type liveResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func NewLiveTrafficLogPolicy() policy.Policy {
	return &liveTrafficLogPolicy{
		notAllowedHeaders: map[string]bool{
			"authorization": true,
		},
	}
}

func (p *liveTrafficLogPolicy) Do(req *policy.Request) (*http.Response, error) {
	rawRequest := req.Raw()
	liveReq := liveRequest{
		Headers: p.header(rawRequest.Header),
		Method:  rawRequest.Method,
		Url:     rawRequest.URL.String(),
		Body:    p.requestBodyString(req),
	}
	if err := req.RewindBody(); err != nil {
		return nil, err
	}
	response, err := req.Next() // Make the request
	liveResp := liveResponse{}
	if err == nil {
		liveResp.Headers = p.header(response.Header)
		liveResp.StatusCode = response.StatusCode
		liveResp.Body = p.responseBodyString(response)
	} else {
		liveResp.Body = err.Error()
	}
	liveTraffic := traffic{
		LiveRequest:  liveReq,
		LiveResponse: liveResp,
	}

	data, marshalErr := json.Marshal(liveTraffic)
	if marshalErr != nil {
		log.Printf("[ERROR] Failed to marshal live traffic: %v", marshalErr)
		return response, err
	}
	log.Printf("[DEBUG] Live traffic: %s", string(data))

	return response, err
}

func (p *liveTrafficLogPolicy) requestBodyString(req *policy.Request) string {
	if req.Raw().Body == nil {
		return ""
	}
	body, err := io.ReadAll(req.Raw().Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %v", err)
		body = []byte(err.Error())
	}
	if err := req.RewindBody(); err != nil {
		log.Printf("[ERROR] Failed to rewind request body: %v", err)
		return ""
	}
	return string(body)
}

func (p *liveTrafficLogPolicy) responseBodyString(resp *http.Response) string {
	body, err := runtime.Payload(resp)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body: %v", err)
		return ""
	}
	return string(body)
}

func (p *liveTrafficLogPolicy) header(input http.Header) map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		if p.notAllowedHeaders[strings.ToLower(k)] {
			output[k] = redactedValue
		} else {
			output[k] = strings.Join(v, ",")
		}
	}
	return output
}
