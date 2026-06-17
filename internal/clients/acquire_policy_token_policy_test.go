package clients

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func newResponseWithBody(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func Test_missingPolicyTokenDetailsFromResponse(t *testing.T) {
	testcases := []struct {
		Name            string
		Body            string
		ExpectError     bool
		ExpectNil       bool
		ExpectedDetails missingPolicyTokenDetails
	}{
		{
			Name: "policy violation requiring a token",
			Body: `{
				"error": {
					"code": "RequestDisallowedByPolicy",
					"additionalInfo": [
						{
							"type": "PolicyViolation",
							"info": {
								"evaluationDetails": {
									"missingPolicyTokenDetails": {
										"shouldDeny": true,
										"isChangeReferenceRequired": false,
										"endpointKind": "CoinFlip"
									}
								}
							}
						}
					]
				}
			}`,
			ExpectedDetails: missingPolicyTokenDetails{
				ShouldDeny:                true,
				IsChangeReferenceRequired: false,
				EndpointKind:              "CoinFlip",
			},
		},
		{
			Name: "policy violation requiring a change reference",
			Body: `{
				"error": {
					"code": "RequestDisallowedByPolicy",
					"additionalInfo": [
						{
							"type": "PolicyViolation",
							"info": {
								"evaluationDetails": {
									"missingPolicyTokenDetails": {
										"shouldDeny": true,
										"isChangeReferenceRequired": true,
										"endpointKind": "CoinFlip"
									}
								}
							}
						}
					]
				}
			}`,
			ExpectedDetails: missingPolicyTokenDetails{
				ShouldDeny:                true,
				IsChangeReferenceRequired: true,
				EndpointKind:              "CoinFlip",
			},
		},
		{
			Name: "error code and type are matched case-insensitively",
			Body: `{
				"error": {
					"code": "requestdisallowedbypolicy",
					"additionalInfo": [
						{
							"type": "policyviolation",
							"info": {
								"evaluationDetails": {
									"missingPolicyTokenDetails": {
										"shouldDeny": true
									}
								}
							}
						}
					]
				}
			}`,
			ExpectedDetails: missingPolicyTokenDetails{
				ShouldDeny: true,
			},
		},
		{
			Name: "the correct additionalInfo entry is selected",
			Body: `{
				"error": {
					"code": "RequestDisallowedByPolicy",
					"additionalInfo": [
						{
							"type": "SomethingElse",
							"info": {}
						},
						{
							"type": "PolicyViolation",
							"info": {
								"evaluationDetails": {
									"missingPolicyTokenDetails": {
										"shouldDeny": true,
										"endpointKind": "CoinFlip"
									}
								}
							}
						}
					]
				}
			}`,
			ExpectedDetails: missingPolicyTokenDetails{
				ShouldDeny:   true,
				EndpointKind: "CoinFlip",
			},
		},
		{
			Name: "policy violation that should not deny is ignored",
			Body: `{
				"error": {
					"code": "RequestDisallowedByPolicy",
					"additionalInfo": [
						{
							"type": "PolicyViolation",
							"info": {
								"evaluationDetails": {
									"missingPolicyTokenDetails": {
										"shouldDeny": false
									}
								}
							}
						}
					]
				}
			}`,
			ExpectNil: true,
		},
		{
			Name: "missingPolicyTokenDetails absent",
			Body: `{
				"error": {
					"code": "RequestDisallowedByPolicy",
					"additionalInfo": [
						{
							"type": "PolicyViolation",
							"info": {
								"evaluationDetails": {}
							}
						}
					]
				}
			}`,
			ExpectNil: true,
		},
		{
			Name: "different error code",
			Body: `{
				"error": {
					"code": "AuthorizationFailed",
					"additionalInfo": [
						{
							"type": "PolicyViolation",
							"info": {
								"evaluationDetails": {
									"missingPolicyTokenDetails": {
										"shouldDeny": true
									}
								}
							}
						}
					]
				}
			}`,
			ExpectNil: true,
		},
		{
			Name: "no policy violation additionalInfo",
			Body: `{
				"error": {
					"code": "RequestDisallowedByPolicy",
					"additionalInfo": [
						{
							"type": "SomethingElse",
							"info": {}
						}
					]
				}
			}`,
			ExpectNil: true,
		},
		{
			Name:      "invalid json",
			Body:      `not json`,
			ExpectNil: true,
		},
		{
			Name:      "empty body",
			Body:      ``,
			ExpectNil: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			resp := newResponseWithBody(http.StatusForbidden, tc.Body)

			details, err := missingPolicyTokenDetailsFromResponse(resp)

			if tc.ExpectError {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.ExpectNil {
				if details != nil {
					t.Fatalf("expected nil details but got %+v", *details)
				}
				return
			}

			if details == nil {
				t.Fatalf("expected details but got nil")
			}
			if *details != tc.ExpectedDetails {
				t.Fatalf("unexpected details: got %+v, want %+v", *details, tc.ExpectedDetails)
			}
		})
	}
}
