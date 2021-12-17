package location

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
	return Normalize(old) == Normalize(new)
}

func LocationStateFunc(location interface{}) string {
	input := location.(string)
	return Normalize(input)
}

func Normalize(input string) string {
	return strings.ReplaceAll(strings.ToLower(input), " ", "")
}
