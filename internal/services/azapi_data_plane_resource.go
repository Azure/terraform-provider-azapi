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

func ResourceAzApiDataPlaneResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzApiDataPlaneResourceCreateUpdate,
		Read:   resourceAzApiDataPlaneResourceRead,
		Update: resourceAzApiDataPlaneResourceCreateUpdate,
		Delete: resourceAzApiDataPlaneResourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return nil, fmt.Errorf("importing is not supported for this resource yet")
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"parent_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
				Default:  true,
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
				// #nosec G104
				d.SetNewComputed("output")
			}
			old, new := d.GetChange("body")
			if utils.NormalizeJson(old) != utils.NormalizeJson(new) {
				// #nosec G104
				d.SetNewComputed("output")
			}
			return nil
		},
	}
}

func resourceAzApiDataPlaneResourceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataPlaneClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !d.IsNewResource() {
		d.Partial(true)
	}

	id, err := parse.NewDataPlaneResourceId(d.Get("name").(string), d.Get("parent_id").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		_, err := client.Get(ctx, id)
		if err == nil {
			return tf.ImportAsExistsError("azapi_data_plane_resource", id.ID())
		}
		if !utils.ResponseErrorWasNotFound(err) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	var body map[string]interface{}
	err = json.Unmarshal([]byte(d.Get("body").(string)), &body)
	if err != nil {
		return err
	}

	for _, id := range d.Get("locks").([]interface{}) {
		locks.ByID(id.(string))
		defer locks.UnlockByID(id.(string))
	}

	_, err = client.CreateOrUpdateThenPoll(ctx, id, body)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	if !d.IsNewResource() {
		d.Partial(false)
	}

	return resourceAzApiDataPlaneResourceRead(d, meta)
}

func resourceAzApiDataPlaneResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataPlaneClient
	ctx, cancel := tf.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataPlaneResourceIDWithResourceType(d.Id(), d.Get("type").(string))
	if err != nil {
		return err
	}

	responseBody, err := client.Get(ctx, id)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			log.Printf("[INFO] Error reading %q - removing from state", id.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %q: %+v", id, err)
	}

	bodyJson := d.Get("body").(string)
	var requestBody interface{}
	err = json.Unmarshal([]byte(bodyJson), &requestBody)
	if err != nil && len(bodyJson) != 0 {
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
	// #nosec G104
	d.Set("body", string(data))

	// #nosec G104
	d.Set("name", id.Name)
	// #nosec G104
	d.Set("parent_id", id.ParentId)
	// #nosec G104
	d.Set("type", fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion))

	// #nosec G104
	d.Set("output", flattenOutput(responseBody, d.Get("response_export_values").([]interface{})))
	return nil
}

func resourceAzApiDataPlaneResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataPlaneClient
	ctx, cancel := tf.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataPlaneResourceIDWithResourceType(d.Id(), d.Get("type").(string))
	if err != nil {
		return err
	}

	for _, id := range d.Get("locks").([]interface{}) {
		locks.ByID(id.(string))
		defer locks.UnlockByID(id.(string))
	}

	_, err = client.DeleteThenPoll(ctx, id)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			return nil
		}
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}
