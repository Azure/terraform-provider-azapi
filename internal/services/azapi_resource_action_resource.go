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

func ResourceResourceAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceResourceActionCreateUpdate,
		Read:   resourceResourceActionRead,
		Update: resourceResourceActionCreateUpdate,
		Delete: resourceResourceActionDelete,

		Importer: tf.DefaultImporter(func(id string) error {
			return fmt.Errorf("`azapi_resource_action` doesn't support import")
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validate.ResourceType,
				RequiredWith:  []string{"resource_id", "action"},
				ConflictsWith: []string{"uri"},
			},

			"resource_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validate.ResourceID,
				RequiredWith:  []string{"type", "action"},
				ConflictsWith: []string{"uri"},
			},

			"action": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				RequiredWith:  []string{"type", "resource_id"},
				ConflictsWith: []string{"uri"},
			},

			"uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"type", "resource_id", "action"},
				RequiredWith:  []string{"api_version"},
			},

			"api_version": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"type", "resource_id", "action"},
				RequiredWith:  []string{"uri"},
			},

			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "POST",
				ValidateFunc: validation.StringInSlice([]string{
					"POST",
					"PATCH",
					"PUT",
					"DELETE",
				}, false),
			},

			"body": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: tf.SuppressJsonOrderingDifference,
				StateFunc:        utils.NormalizeJson,
			},

			"locks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"response_export_values": {
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
			if d.HasChange("response_export_values") || d.HasChange("action") || d.HasChange("body") {
				// #nosec G104
				d.SetNewComputed("output")
			}
			return nil
		},
	}
}

func resourceResourceActionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !d.IsNewResource() {
		d.Partial(true)
	}

	actionName, actionOk := d.GetOk("action")
	uri, uriOk := d.GetOk("uri")
	apiVersion, apiVersionOk := d.GetOk("api_version")
	resourceId, resourceIdOk := d.GetOk("resource_id")
	typeName, typeOk := d.GetOk("type")

	if !actionOk {
		actionName = ""
	}

	var id parse.ResourceId
	if actionOk && resourceIdOk && typeOk {
		var err error
		id, err = parse.ResourceIDWithResourceType(resourceId.(string), typeName.(string))
		if err != nil {
			return err
		}
	}

	if uriOk && apiVersionOk {
		id = parse.ResourceId{
			AzureResourceId: uri.(string),
			ApiVersion:      apiVersion.(string),
		}
	}

	method := d.Get("method").(string)
	body := d.Get("body").(string)
	var requestBody interface{}
	if len(body) != 0 {
		var err error
		err = json.Unmarshal([]byte(body), &requestBody)
		if err != nil {
			return fmt.Errorf("unmarshalling `body`: %+v", err)
		}
	}

	log.Printf("[INFO] request body: %v\n", body)

	for _, id := range d.Get("locks").([]interface{}) {
		locks.ByID(id.(string))
		defer locks.UnlockByID(id.(string))
	}

	responseBody, err := client.Action(ctx, id.AzureResourceId, actionName.(string), id.ApiVersion, method, requestBody)
	if err != nil {
		return fmt.Errorf("performing action %s of %q: %+v", actionName, id, err)
	}

	resourceIdplusApi := id.ID()
	if actionName != "" {
		resourceIdplusApi = fmt.Sprintf("%s/%s", id.ID(), actionName)
	}
	d.SetId(resourceIdplusApi)
	// #nosec G104
	d.Set("output", flattenOutput(responseBody, d.Get("response_export_values").([]interface{})))

	if !d.IsNewResource() {
		d.Partial(false)
	}

	return resourceResourceActionRead(d, meta)
}

func resourceResourceActionRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceResourceActionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
