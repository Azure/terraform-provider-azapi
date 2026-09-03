package services

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func Test_withApiVersionFallback(t *testing.T) {
	statusErr := func(code int) error {
		return &azcore.ResponseError{StatusCode: code}
	}

	testcases := []struct {
		name                 string
		operation            func(calls *[]string) func(string) (interface{}, error)
		explicitApiVersion   string
		candidateApiVersions []string
		expectBody           interface{}
		expectErr            bool
		expectCalls          []string
	}{
		{
			name: "explicit api version is used and not retried",
			operation: func(calls *[]string) func(string) (interface{}, error) {
				return func(v string) (interface{}, error) {
					*calls = append(*calls, v)
					return nil, statusErr(http.StatusBadRequest)
				}
			},
			explicitApiVersion:   "2015-01-01",
			candidateApiVersions: []string{"2018-01-01", "2019-01-01", "2020-01-01"},
			expectBody:           nil,
			expectErr:            true,
			expectCalls:          []string{"2015-01-01"},
		},
		{
			name: "succeeds on first candidate (latest)",
			operation: func(calls *[]string) func(string) (interface{}, error) {
				return func(v string) (interface{}, error) {
					*calls = append(*calls, v)
					return "ok", nil
				}
			},
			candidateApiVersions: []string{"2018-01-01", "2019-01-01", "2020-01-01"},
			expectBody:           "ok",
			expectErr:            false,
			expectCalls:          []string{"2020-01-01"},
		},
		{
			name: "retries from latest backwards on 404 then succeeds",
			operation: func(calls *[]string) func(string) (interface{}, error) {
				return func(v string) (interface{}, error) {
					*calls = append(*calls, v)
					if v == "2019-01-01" {
						return "ok", nil
					}
					return nil, statusErr(http.StatusNotFound)
				}
			},
			candidateApiVersions: []string{"2018-01-01", "2019-01-01", "2020-01-01"},
			expectBody:           "ok",
			expectErr:            false,
			expectCalls:          []string{"2020-01-01", "2019-01-01"},
		},
		{
			name: "retries at most 3 times",
			operation: func(calls *[]string) func(string) (interface{}, error) {
				return func(v string) (interface{}, error) {
					*calls = append(*calls, v)
					return nil, statusErr(http.StatusBadRequest)
				}
			},
			candidateApiVersions: []string{"2017-01-01", "2018-01-01", "2019-01-01", "2020-01-01"},
			expectBody:           nil,
			expectErr:            true,
			expectCalls:          []string{"2020-01-01", "2019-01-01", "2018-01-01"},
		},
		{
			name: "does not retry on non 400/404 error",
			operation: func(calls *[]string) func(string) (interface{}, error) {
				return func(v string) (interface{}, error) {
					*calls = append(*calls, v)
					return nil, statusErr(http.StatusForbidden)
				}
			},
			candidateApiVersions: []string{"2018-01-01", "2019-01-01", "2020-01-01"},
			expectBody:           nil,
			expectErr:            true,
			expectCalls:          []string{"2020-01-01"},
		},
		{
			name: "stops retrying once a non 400/404 error occurs",
			operation: func(calls *[]string) func(string) (interface{}, error) {
				return func(v string) (interface{}, error) {
					*calls = append(*calls, v)
					if v == "2020-01-01" {
						return nil, statusErr(http.StatusNotFound)
					}
					return nil, errors.New("boom")
				}
			},
			candidateApiVersions: []string{"2018-01-01", "2019-01-01", "2020-01-01"},
			expectBody:           nil,
			expectErr:            true,
			expectCalls:          []string{"2020-01-01", "2019-01-01"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var calls []string
			body, err := withApiVersionFallback(context.Background(), tc.operation(&calls), tc.explicitApiVersion, tc.candidateApiVersions)
			if tc.expectErr && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.expectErr && err != nil {
				t.Fatalf("expected no error but got: %+v", err)
			}
			if !reflect.DeepEqual(body, tc.expectBody) {
				t.Fatalf("expected body %v but got %v", tc.expectBody, body)
			}
			if !reflect.DeepEqual(calls, tc.expectCalls) {
				t.Fatalf("expected calls %v but got %v", tc.expectCalls, calls)
			}
		})
	}
}
