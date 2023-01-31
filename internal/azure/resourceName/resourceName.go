package resourceName

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func SchemaResourceNameOC() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	}
}

func SchemaResourceName() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"default_naming_prefix", "default_naming_suffix"},
	}
}

func SchemaResourceNamePrefix() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"default_name"},
	}
}

func SchemaResourceNameSuffix() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"default_name"},
	}
}

func SchemaResourceNameRemovingSpecialCharacters() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
}
