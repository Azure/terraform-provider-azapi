package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/clients"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/services/validate"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/tf"
	"github.com/Azure/terraform-provider-azurerm-restapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAzureGenericPatchResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureGenericPatchResourceCreateUpdate,
		Read:   resourceAzureGenericPatchResourceRead,
		Update: resourceAzureGenericPatchResourceCreateUpdate,
		Delete: resourceAzureGenericPatchResourceDelete,

		Importer: tf.DefaultImporter(func(id string) error {
			return fmt.Errorf("`azurerm-restapi_patch_resource` doesn't support import")
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				RequiredWith:  []string{"parent_id"},
				ConflictsWith: []string{"resource_id"},
				AtLeastOneOf:  []string{"name", "resource_id"},
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				//ValidateFunc:  validate.AzureResourceID,
				RequiredWith:  []string{"name"},
				ConflictsWith: []string{"resource_id"},
			},

			"resource_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  validate.AzureResourceID,
				ConflictsWith: []string{"name", "parent_id"},
				AtLeastOneOf:  []string{"name", "resource_id"},
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

			if name := d.Get("name").(string); len(name) != 0 {
				id, err := parse.BuildResourceID(d.Get("name").(string), d.Get("parent_id").(string), d.Get("type").(string))
				if err != nil && len(id.ParentId) > 0 {
					return err
				}
			}

			return nil
		},
	}
}

func resourceAzureGenericPatchResourceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var id parse.ResourceId
	if name := d.Get("name").(string); len(name) != 0 {
		buildId, err := parse.BuildResourceID(d.Get("name").(string), d.Get("parent_id").(string), d.Get("type").(string))
		if err != nil {
			return err
		}
		id = buildId
	} else {
		buildId, err := parse.NewResourceID(d.Get("resource_id").(string), d.Get("type").(string))
		if err != nil {
			return err
		}
		id = buildId
	}

	existing, _, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}
	if len(utils.GetId(existing)) == 0 {
		return fmt.Errorf("patch target does not exist %s", id)
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
	_, _, err = client.CreateUpdate(ctx, id.AzureResourceId, id.ApiVersion, requestBody, http.MethodPut)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAzureGenericPatchResourceRead(d, meta)
}

func resourceAzureGenericPatchResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceID(d.Id())
	if err != nil {
		return err
	}

	responseBody, response, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if response.StatusCode == http.StatusNotFound {
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
	data, err := json.Marshal(utils.GetUpdatedJson(requestBody, responseBody))
	if err != nil {
		return err
	}
	d.Set("body", string(data))

	paths := d.Get("response_export_values").([]interface{})
	var output interface{}
	if len(paths) != 0 {
		output = make(map[string]interface{})
		for _, path := range paths {
			part := utils.ExtractObject(responseBody, path.(string))
			if part == nil {
				continue
			}
			output = utils.GetMergedJson(output, part)
		}
	}
	if output == nil {
		output = make(map[string]interface{})
	}
	outputJson, _ := json.Marshal(output)
	d.Set("output", string(outputJson))
	return nil
}

func resourceAzureGenericPatchResourceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
