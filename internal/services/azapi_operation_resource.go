package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAzApiOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzApiOperationCreateUpdate,
		Read:   resourceAzApiOperationRead,
		Update: resourceAzApiOperationCreateUpdate,
		Delete: resourceAzApiOperationDelete,

		Importer: tf.DefaultImporter(func(id string) error {
			return fmt.Errorf("`azapi_operation` doesn't support import")
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
				ValidateFunc: validate.ResourceType,
			},

			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AzureResourceID,
			},

			"operation": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "POST",
				ValidateFunc: validation.StringInSlice([]string{
					"POST",
					"PATCH",
					"GET",
					"PUT",
					"DELETE",
					"CONNECT",
					"HEAD",
					"OPTIONS",
					"TRACE",
				},
					false,
				),
			},

			"body": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: tf.SuppressJsonOrderingDifference,
			},

			"response_export_values": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
			return nil
		},
	}
}

func resourceAzApiOperationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NewResourceID(d.Get("resource_id").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	operationName := d.Get("operation").(string)
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
	responseBody, err := client.Action(ctx, id.AzureResourceId, operationName, id.ApiVersion, method, requestBody)
	if err != nil {
		return fmt.Errorf("performing action %s of %q: %+v", operationName, id, err)
	}

	d.SetId(id.ID())
	d.Set("output", flattenOutput(responseBody, d.Get("response_export_values").([]interface{})))

	return resourceAzApiOperationRead(d, meta)
}

func resourceAzApiOperationRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAzApiOperationDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
