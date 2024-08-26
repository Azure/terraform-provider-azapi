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
