package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func AzApiDataSource() *schema.Resource {
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
				ExactlyOneOf: []string{"name", "resource_id"},
			},

			"parent_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ResourceID,
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
		parentId := d.Get("parent_id").(string)
		resourceType := d.Get("type").(string)
		if parentId == "" && strings.HasPrefix(strings.ToUpper(resourceType), strings.ToUpper(arm.ResourceGroupResourceType.String())) {
			parentId = fmt.Sprintf("/subscriptions/%s", meta.(*clients.Client).Account.GetSubscriptionId())
		}

		buildId, err := parse.NewResourceID(name, parentId, resourceType)
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

	var responseBody interface{}
	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			return fmt.Errorf("not found %q: %+v", id, err)
		}
		return fmt.Errorf("reading %q: %+v", id, err)
	}

	// we retry only if retry if responseBody is nil
	// this helps to retry for few attributes of an Azure resource to be updated by Azure policies etc.,
	// outside the scope of Terraform
	if responseBody == nil {
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			responseBody, err = client.Get(ctx, id.AzureResourceId, id.ApiVersion)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if responseBody != nil {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("data provider %s doesn't exist yet, retrying", id.AzureResourceId))
		})
		if err != nil {
			return fmt.Errorf("error fetching the data provider inside retry function: %s", err)
		}
	}

	d.SetId(id.ID())
	// #nosec G104
	d.Set("name", id.Name)
	// #nosec G104
	d.Set("parent_id", id.ParentId)
	if bodyMap, ok := responseBody.(map[string]interface{}); ok {
		// #nosec G104
		d.Set("tags", tags.FlattenTags(bodyMap["tags"]))
		// #nosec G104
		d.Set("location", bodyMap["location"])
		// #nosec G104
		d.Set("identity", identity.FlattenIdentity(bodyMap["identity"]))
	}
	// #nosec G104
	d.Set("output", flattenOutput(responseBody, d.Get("response_export_values").([]interface{})))
	return nil
}
