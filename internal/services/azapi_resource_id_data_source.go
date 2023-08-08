package services

import (
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIdDataSource() *schema.Resource {
	return &schema.Resource{
		Read: resourceIdDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ResourceType,
			},

			"parent_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ResourceID,
				RequiredWith: []string{"name"},
				ExactlyOneOf: []string{"resource_id", "parent_id"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"parent_id"},
			},

			"resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ResourceID,
				ExactlyOneOf: []string{"resource_id", "parent_id"},
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subscription_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"provider_namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"parts": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceIdDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	_, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
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

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("parent_id", id.ParentId)
	d.Set("resource_id", id.AzureResourceId)

	armId, err := arm.ParseResourceID(id.AzureResourceId)
	if id.AzureResourceId == "/" {
		armId, err = &arm.ResourceID{
			ResourceType: arm.TenantResourceType,
		}, nil
	}
	if err != nil {
		return err
	}

	d.Set("resource_group_name", armId.ResourceGroupName)
	d.Set("subscription_id", armId.SubscriptionID)
	d.Set("provider_namespace", armId.ResourceType.Namespace)

	path := id.AzureResourceId
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	components := strings.Split(path, "/")
	parts := make(map[string]string)
	for i := 0; i < len(components)-1; i += 2 {
		parts[components[i]] = components[i+1]
	}
	d.Set("parts", parts)

	return nil
}
