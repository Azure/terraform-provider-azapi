package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceResourceActionDataSource() *schema.Resource {
	return &schema.Resource{
		Read: resourceResourceActionDataSourceRead,

		Importer: tf.DefaultImporter(func(id string) error {
			return fmt.Errorf("`azapi_resource_action` doesn't support import")
		}),

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.ResourceID,
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ResourceType,
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
					"GET",
				}, false),
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
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"output": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceResourceActionDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
	responseBody, err := client.Action(ctx, id.AzureResourceId, actionName, id.ApiVersion, method, requestBody)
	if err != nil {
		return fmt.Errorf("performing action %s of %q: %+v", actionName, id, err)
	}

	resourceId := id.ID()
	if actionName != "" {
		resourceId = fmt.Sprintf("%s/%s", id.ID(), actionName)
	}
	d.SetId(resourceId)
	d.Set("output", flattenOutput(responseBody, d.Get("response_export_values").([]interface{})))

	return resourceResourceActionRead(d, meta)
}
