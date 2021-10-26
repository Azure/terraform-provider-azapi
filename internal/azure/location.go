package azure

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SchemaLocation() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Computed:         true,
		StateFunc:        LocationStateFunc,
		DiffSuppressFunc: LocationDiffSuppressFunc,
	}
}

func SchemaLocationDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

func LocationDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return LocationNormalize(old) == LocationNormalize(new)
}

func LocationStateFunc(location interface{}) string {
	input := location.(string)
	return LocationNormalize(input)
}

func LocationNormalize(input string) string {
	return strings.ReplaceAll(strings.ToLower(input), " ", "")
}

func ExpandLocation(location string) interface{} {
	body := make(map[string]interface{}, 0)
	body["location"] = LocationNormalize(location)
	return body
}

func FlattenLocation(body interface{}) interface{} {
	if body != nil {
		if bodyMap, ok := body.(map[string]interface{}); ok && bodyMap["location"] != nil {
			return bodyMap["location"].(string)
		}
	}
	return nil
}
