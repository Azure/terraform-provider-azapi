package services

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func Test_countCoveredLeaves(t *testing.T) {
	testcases := []struct {
		Name          string
		Body          string
		Response      string
		ExpectTotal   int
		ExpectCovered int
	}{
		{
			Name:          "empty get response drops every configured leaf",
			Body:          `{"properties":{"MY_SETTING":"value1","OTHER":"value2"}}`,
			Response:      `{"properties":{},"name":"appsettings","kind":null}`,
			ExpectTotal:   2,
			ExpectCovered: 0,
		},
		{
			Name:          "list response covers every configured leaf",
			Body:          `{"properties":{"MY_SETTING":"value1","OTHER":"value2"}}`,
			Response:      `{"properties":{"MY_SETTING":"value1","OTHER":"value2"},"name":"appsettings"}`,
			ExpectTotal:   2,
			ExpectCovered: 2,
		},
		{
			Name:          "partial coverage still counts as covered",
			Body:          `{"properties":{"MY_SETTING":"value1","OTHER":"value2"}}`,
			Response:      `{"properties":{"MY_SETTING":"value1"}}`,
			ExpectTotal:   2,
			ExpectCovered: 1,
		},
		{
			Name:          "nil response covers nothing",
			Body:          `{"properties":{"MY_SETTING":"value1"}}`,
			Response:      `null`,
			ExpectTotal:   1,
			ExpectCovered: 0,
		},
		{
			Name:          "empty body has no leaves",
			Body:          `{}`,
			Response:      `{"properties":{"MY_SETTING":"value1"}}`,
			ExpectTotal:   0,
			ExpectCovered: 0,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			var body, response interface{}
			_ = json.Unmarshal([]byte(testcase.Body), &body)
			_ = json.Unmarshal([]byte(testcase.Response), &response)

			total, covered := countCoveredLeaves(body, response)
			if total != testcase.ExpectTotal || covered != testcase.ExpectCovered {
				t.Fatalf("expected total=%d covered=%d but got total=%d covered=%d", testcase.ExpectTotal, testcase.ExpectCovered, total, covered)
			}
		})
	}
}

func Test_applyResponsePath(t *testing.T) {
	var response interface{}
	_ = json.Unmarshal([]byte(`{"value":{"properties":{"MY_SETTING":"value1"}}}`), &response)

	var diags diag.Diagnostics
	result := applyResponsePath(response, "value", &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	var expected interface{}
	_ = json.Unmarshal([]byte(`{"properties":{"MY_SETTING":"value1"}}`), &expected)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %v but got %v", expected, result)
	}

	result = applyResponsePath(response, "", &diags)
	if !reflect.DeepEqual(result, response) {
		t.Fatalf("empty path should return the response unchanged")
	}
}
