package services

import (
	"fmt"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAzApiDataSource() *schema.Resource {
	return &schema.Resource{
		Read: resourceAzApiDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				RequiredWith: []string{"parent_id"},
				ExactlyOneOf: []string{"name", "resource_id"},
			},

			"parent_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ResourceID,
				RequiredWith: []string{"name"},
			},

			"resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ResourceID,
				ExactlyOneOf: []string{"name", "resource_id"},
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ResourceType,
			},

			"response_export_values": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"location": location.SchemaLocationDataSource(),

			"identity": identity.SchemaIdentityDataSource(),

			"output": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaTagsDataSource(),
		},
	}
}

func resourceAzApiDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var id parse.ResourceId
	if name := d.Get("name").(string); len(name) != 0 {
		buildId, err := parse.NewResourceID(d.Get("name").(string), d.Get("parent_id").(string), d.Get("type").(string))
		if err != nil {
			return err
		}
		id = buildId
	} else {
		buildId, err := parse.ResourceIDWithResourceType(d.Get("resource_id").(string), d.Get("type").(string))
		if err != nil {
			return err
		}
		id = buildId
	}

	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			return fmt.Errorf("not found %q: %+v", id, err)
		}
		return fmt.Errorf("reading %q: %+v", id, err)
	}
	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("parent_id", id.ParentId)
	if bodyMap, ok := responseBody.(map[string]interface{}); ok {
		d.Set("tags", tags.FlattenTags(bodyMap["tags"]))
		d.Set("location", bodyMap["location"])
		d.Set("identity", identity.FlattenIdentity(bodyMap["identity"]))
	}
	d.Set("output", flattenOutput(responseBody, d.Get("response_export_values").([]interface{})))
	return nil
}
