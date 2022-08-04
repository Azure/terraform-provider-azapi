package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/locks"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAzApiUpdateResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzApiUpdateResourceCreateUpdate,
		Read:   resourceAzApiUpdateResourceRead,
		Update: resourceAzApiUpdateResourceCreateUpdate,
		Delete: resourceAzApiUpdateResourceDelete,

		Importer: tf.DefaultImporter(func(id string) error {
			return fmt.Errorf("`azapi_update_resource` doesn't support import")
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				RequiredWith: []string{"parent_id"},
				ExactlyOneOf: []string{"name", "resource_id"},
			},

			"parent_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validate.ResourceID,
				RequiredWith: []string{"name"},
			},

			"resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validate.ResourceID,
				ExactlyOneOf: []string{"name", "resource_id"},
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ResourceType,
			},

			"body": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: tf.SuppressJsonOrderingDifference,
			},

			"ignore_casing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ignore_missing_property": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"response_export_values": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"locks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"output": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			if d.HasChange("response_export_values") {
				d.SetNewComputed("output")
			}
			old, new := d.GetChange("body")
			if utils.NormalizeJson(old) != utils.NormalizeJson(new) {
				d.SetNewComputed("output")
			}

			if name := d.Get("name").(string); len(name) != 0 {
				parentId := d.Get("parent_id").(string)
				resourceType := d.Get("type").(string)

				// verify parent_id when it's known
				if len(parentId) > 0 {
					_, err := parse.NewResourceID(name, parentId, resourceType)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}

func resourceAzApiUpdateResourceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
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

	existing, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}
	if utils.GetId(existing) == nil {
		return fmt.Errorf("update target does not exist %s", id)
	}

	var requestBody interface{}
	err = json.Unmarshal([]byte(d.Get("body").(string)), &requestBody)
	if err != nil {
		return err
	}

	requestBody = utils.GetMergedJson(existing, requestBody)
	if id.ResourceDef != nil {
		requestBody = (*id.ResourceDef).GetWriteOnly(requestBody)
	}
	j, _ := json.Marshal(requestBody)
	log.Printf("[INFO] request body: %v\n", string(j))
	for _, id := range d.Get("locks").([]interface{}) {
		locks.ByID(id.(string))
		defer locks.UnlockByID(id.(string))
	}
	_, err = client.CreateOrUpdate(ctx, id.AzureResourceId, id.ApiVersion, requestBody)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAzApiUpdateResourceRead(d, meta)
}

func resourceAzApiUpdateResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var id parse.ResourceId
	var err error
	if resourceType := d.Get("type").(string); len(resourceType) != 0 {
		id, err = parse.ResourceIDWithResourceType(d.Id(), resourceType)
	} else {
		id, err = parse.ResourceIDWithApiVersion(d.Id())
	}
	if err != nil {
		return err
	}

	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			log.Printf("[INFO] Error reading %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("parent_id", id.ParentId)
	d.Set("resource_id", id.AzureResourceId)
	d.Set("type", fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion))

	bodyJson := d.Get("body").(string)
	var requestBody interface{}
	err = json.Unmarshal([]byte(bodyJson), &requestBody)
	if err != nil {
		return err
	}
	option := utils.UpdateJsonOption{
		IgnoreCasing:          d.Get("ignore_casing").(bool),
		IgnoreMissingProperty: d.Get("ignore_missing_property").(bool),
	}
	data, err := json.Marshal(utils.GetUpdatedJson(requestBody, responseBody, option))
	if err != nil {
		return err
	}
	d.Set("body", string(data))
	d.Set("output", flattenOutput(responseBody, d.Get("response_export_values").([]interface{})))
	return nil
}

func resourceAzApiUpdateResourceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
