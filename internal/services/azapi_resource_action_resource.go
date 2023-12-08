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

func ResourceAction() *schema.Resource {
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ResourceType,
			},

			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ResourceID,
			},

			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
					"GET",
					"HEAD",
				}, false),
			},

			"body": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: tf.SuppressJsonOrderingDifference,
				StateFunc:        utils.NormalizeJson,
			},

			"when": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "apply",
				Description:  "When to perform the action, value must be one of: 'apply', 'destroy'. Default is 'apply'.",
				ValidateFunc: validation.StringInSlice([]string{"apply", "destroy"}, false),
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

	id, err := parse.ResourceIDWithResourceType(d.Get("resource_id").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	actionName := d.Get("action").(string)
	method := d.Get("method").(string)
	body := d.Get("body").(string)
	var requestBody interface{}
	if len(body) != 0 {
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

	var responseBody interface{} = "{}"
	if d.Get("when").(string) == "apply" {
		responseBody, err = client.Action(ctx, id.AzureResourceId, actionName, id.ApiVersion, method, requestBody)
		if err != nil {
			return fmt.Errorf("performing action %s of %q: %+v", actionName, id, err)
		}
	}

	resourceId := id.ID()
	if actionName != "" {
		resourceId = fmt.Sprintf("%s/%s", id.ID(), actionName)
	}
	d.SetId(resourceId)

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
	if when, _ := d.GetChange("when"); when.(string) != "destroy" {
		return nil
	}

	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceIDWithResourceType(d.Get("resource_id").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	actionName, _ := d.GetChange("action")
	method, _ := d.GetChange("method")
	body, _ := d.GetChange("body")
	var requestBody interface{}
	if len(body.(string)) != 0 {
		err := json.Unmarshal([]byte(body.(string)), &requestBody)
		if err != nil {
			return fmt.Errorf("unmarshalling `body`: %+v", err)
		}
	}

	if _, err := client.Action(ctx, id.AzureResourceId, actionName.(string), id.ApiVersion, method.(string), requestBody); err != nil {
		return fmt.Errorf("performing action %s of %q: %+v", actionName, id, err)
	}
	return nil
}
