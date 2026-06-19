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

func Test_policyTokenFromResponse(t *testing.T) {
	testcases := []struct {
		Name                  string
		Body                  string
		ExpectError           bool
		ExpectedToken         string
		ExpectedErrorContains string
	}{
		{
			Name: "succeeded with token",
			Body: `{
				"result": "Succeeded",
				"token": "the-policy-token"
			}`,
			ExpectedToken: "the-policy-token",
		},
		{
			Name: "result is matched case-insensitively",
			Body: `{
				"result": "succeeded",
				"token": "the-policy-token"
			}`,
			ExpectedToken: "the-policy-token",
		},
		{
			Name: "failed result returns an error with the result messages",
			Body: `{
				"result": "Failed",
				"results": [
					{
						"policyInfo": {
							"policyDefinitionId": "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/exampleDenyAction",
							"policyDefinitionName": "exampleDenyAction",
							"policyDefinitionDisplayName": "exampleDenyAction",
							"policyDefinitionVersion": "1.0.0",
							"policyDefinitionEffect": "deny",
							"policyAssignmentId": "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/exampleDenyAction",
							"policyAssignmentName": "exampleDenyAction",
							"policyAssignmentDisplayName": "exampleDenyAction",
							"policyAssignmentScope": "/subscriptions/00000000-0000-0000-0000-000000000000",
							"policyExemptionIds": []
						},
						"result": "Failed",
						"message": "Create validation failed with resource id: '/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/exampleRg/providers/Microsoft.Compute/virtualMachineScaleSets/exampleTest' and validator id: '/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ChangeSafety/validators/metricsValidator/versions/0.0.1-beta'. Error code: 'AuthorizationFailed'."
					}
				]
			}`,
			ExpectError:           true,
			ExpectedErrorContains: "Create validation failed with resource id",
		},
		{
			Name: "failed result with multiple messages are joined",
			Body: `{
				"result": "Failed",
				"results": [
					{ "result": "Failed", "message": "first failure" },
					{ "result": "Failed", "message": "second failure" }
				]
			}`,
			ExpectError:           true,
			ExpectedErrorContains: "first failure; second failure",
		},
		{
			Name: "failed result without any messages still returns an error",
			Body: `{
				"result": "Failed"
			}`,
			ExpectError:           true,
			ExpectedErrorContains: "without any details",
		},
		{
			Name: "succeeded without a token is treated as no token",
			Body: `{
				"result": "Succeeded"
			}`,
			ExpectedToken: "",
		},
		{
			Name:        "invalid json returns an error",
			Body:        `not json`,
			ExpectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			resp := newResponseWithBody(http.StatusOK, tc.Body)

			token, err := policyTokenFromResponse(resp)

			if tc.ExpectError {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				if tc.ExpectedErrorContains != "" && !strings.Contains(err.Error(), tc.ExpectedErrorContains) {
					t.Fatalf("expected error to contain %q, got %q", tc.ExpectedErrorContains, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if token != tc.ExpectedToken {
				t.Fatalf("unexpected token: got %q, want %q", token, tc.ExpectedToken)
			}
		})
	}
}
