package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/clients"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/parse"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/tf"
	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

func ResourceAzureGenericPatchResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureGenericPatchResourceCreateUpdate,
		Read:   resourceAzureGenericPatchResourceRead,
		Update: resourceAzureGenericPatchResourceCreateUpdate,
		Delete: resourceAzureGenericPatchResourceDelete,

		Importer: tf.DefaultImporter(func(id string) error {
			_, err := parse.ResourceID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"api_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"body": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: tf.SuppressJsonOrderingDifference,
			},

			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "PUT",
				ValidateFunc: validation.StringInSlice([]string{
					http.MethodPost,
					http.MethodPut,
					// http.MethodPatch, not supported yet
				}, false),
			},

			"response_body": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAzureGenericPatchResourceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewResourceID(d.Get("url").(string), d.Get("api_version").(string))

	existing, _, err := client.Get(ctx, id.Url, id.ApiVersion)
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
	requestBody = utils.GetIgnoredJson(requestBody, getUnsupportedProperties())
	_, _, err = client.CreateUpdate(ctx, id.Url, id.ApiVersion, requestBody, d.Get("method").(string))
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

	responseBody, response, err := client.Get(ctx, id.Url, id.ApiVersion)
	if err != nil {
		if response.StatusCode == http.StatusNotFound {
			log.Printf("[INFO] Error reading %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	d.Set("url", id.Url)
	d.Set("api_version", id.ApiVersion)

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

	responseBodyJson, err := json.Marshal(responseBody)
	d.Set("response_body", string(responseBodyJson))
	return nil
}

func resourceAzureGenericPatchResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceID(d.Id())
	if err != nil {
		return err
	}

	existing, _, err := client.Get(ctx, id.Url, id.ApiVersion)
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
	requestBody = utils.GetIgnoredJson(requestBody, getUnsupportedProperties())
	_, _, err = client.CreateUpdate(ctx, id.Url, id.ApiVersion, requestBody, d.Get("method").(string))
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	return nil
}

func getUnsupportedProperties() []string {
	return []string{"provisioningState"}
}
