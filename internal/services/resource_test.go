package services

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_FlattenOutputJMES(t *testing.T) {
	testcases := []struct {
		ResponseBody string
		Paths        map[string]string
		ExpectJson   string
	}{
		{
			ResponseBody: `
{
  "values": [
    {
      "id": "1",
      "name": "test1",
      "properties": {
        "status": "active"
      }
    },
    {
      "id": "2",
      "name": "test2",
      "properties": {
        "status": "inactive"
      }
    }
  ]
}
`,
			Paths: map[string]string{
				"names":    "values[*].name",
				"statuses": "values[*].properties.status",
			},
			ExpectJson: `
{
  "names": ["test1", "test2"],
  "statuses": ["active", "inactive"]
}
`,
		},
		{
			ResponseBody: `
{
  "values": [
    {
      "id": "1",
      "name": "test1",
      "properties": {
        "status": "active"
      }
    },
    {
      "id": "2",
      "name": "test2",
      "properties": {
        "status": "inactive"
      }
    }
  ]
}
`,
			Paths: map[string]string{
				"names":       "values[*].name",
				"nonexistent": "values[*].nonexistent",
			},
			ExpectJson: `
{
  "names": ["test1", "test2"],
  "nonexistent": []
}
`,
		},
		{
			ResponseBody: `
{
  "values": [
    {
      "id": "1",
      "name": "test1",
      "properties": {
        "status": "active"
      }
    },
    {
      "id": "2",
      "name": "test2",
      "properties": {
        "status": "inactive"
      }
    }
  ]
}
`,
			Paths: map[string]string{
				"names": "values[*].name",
			},
			ExpectJson: `
{
  "names": ["test1", "test2"]
}
`,
		},
	}

	for _, testcase := range testcases {
		var responseBody, expected interface{}
		_ = json.Unmarshal([]byte(testcase.ResponseBody), &responseBody)
		_ = json.Unmarshal([]byte(testcase.ExpectJson), &expected)

		resultData := flattenOutputJMES(responseBody, testcase.Paths)
		resultJson, _ := dynamic.ToJSON(resultData.(types.Dynamic))

		var result interface{}
		_ = json.Unmarshal(resultJson, &result)

		if !reflect.DeepEqual(result, expected) {
			expectedJson, _ := json.Marshal(expected)
			resultJson, _ = json.Marshal(result)
			t.Fatalf("Expected %s but got %s", string(expectedJson), string(resultJson))
		}
	}
}

func Test_hasApiVersionParameter(t *testing.T) {
	testcases := []struct {
		Name   string
		ID     string
		Expect bool
	}{
		{
			Name:   "azure resource id with api-version",
			ID:     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/virtualNetworks/example-vnet?api-version=2023-11-01",
			Expect: true,
		},
		{
			Name:   "azure resource id without query",
			ID:     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/virtualNetworks/example-vnet",
			Expect: false,
		},
		{
			Name:   "api-version among multiple query parameters",
			ID:     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg?foo=bar&api-version=2023-11-01",
			Expect: true,
		},
		{
			Name:   "api-version is case insensitive",
			ID:     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg?Api-Version=2023-11-01",
			Expect: true,
		},
		{
			Name:   "empty api-version value",
			ID:     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg?api-version=",
			Expect: false,
		},
		{
			Name:   "query without api-version",
			ID:     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg?foo=bar",
			Expect: false,
		},
		{
			Name:   "empty id",
			ID:     "",
			Expect: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			if got := hasApiVersionParameter(testcase.ID); got != testcase.Expect {
				t.Fatalf("expected %v but got %v for ID %q", testcase.Expect, got, testcase.ID)
			}
		})
	}
}
