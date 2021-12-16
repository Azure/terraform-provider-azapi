package tags

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SchemaTagsOC() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeMap,
		Optional:     true,
		Computed:     true,
		ValidateFunc: ValidateTags,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func SchemaTags() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeMap,
		Optional:     true,
		ValidateFunc: ValidateTags,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func SchemaTagsDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func ValidateTags(v interface{}, _ string) (warnings []string, errors []error) {
	tagsMap := v.(map[string]interface{})

	if len(tagsMap) > 50 {
		errors = append(errors, fmt.Errorf("a maximum of 50 tags can be applied to each ARM resource"))
	}

	for k, v := range tagsMap {
		if len(k) > 512 {
			errors = append(errors, fmt.Errorf("the maximum length for a tag key is 512 characters: %q is %d characters", k, len(k)))
		}

		value, err := TagValueToString(v)
		if err != nil {
			errors = append(errors, err)
		} else if len(value) > 256 {
			errors = append(errors, fmt.Errorf("the maximum length for a tag value is 256 characters: the value for %q is %d characters", k, len(value)))
		}
	}

	return warnings, errors
}

func TagValueToString(v interface{}) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil
	case int:
		return fmt.Sprintf("%d", value), nil
	default:
		return "", fmt.Errorf("unknown tag type %T in tag value", value)
	}
}

func ExpandTags(tagsMap map[string]interface{}) map[string]string {
	output := make(map[string]string, len(tagsMap))

	for i, v := range tagsMap {
		// Validate should have ignored this error already
		value, _ := TagValueToString(v)
		output[i] = value
	}

	return output
}

func FlattenTags(input interface{}) map[string]interface{} {
	if input != nil {
		return input.(map[string]interface{})
	}
	return nil
}
