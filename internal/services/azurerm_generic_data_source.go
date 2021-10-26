package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/azure"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/clients"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/parse"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/tf"
	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

func ResourceAzureGenericDataSource() *schema.Resource {
	return &schema.Resource{
		Read: resourceAzureGenericDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
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

			"paths": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"location": azure.SchemaLocationDataSource(),

			"identity": azure.SchemaIdentityDataSource(),

			"output": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": azure.SchemaTagsDataSource(),
		},
	}
}

func resourceAzureGenericDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewResourceID(d.Get("url").(string), d.Get("api_version").(string))

	responseBody, response, err := client.Get(ctx, id.Url, id.ApiVersion)
	if err != nil {
		if response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("not found %q: %+v", id, err)
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}
	d.SetId(id.ID())
	d.Set("tags", azure.FlattenTags(responseBody))
	d.Set("location", azure.FlattenLocation(responseBody))
	d.Set("identity", azure.FlattenIdentity(responseBody))

	paths := d.Get("paths").([]interface{})
	var output interface{}
	if len(paths) != 0 {
		output = make(map[string]interface{}, 0)
		for _, path := range paths {
			part := utils.ExtractObject(responseBody, path.(string))
			if part == nil {
				continue
			}
			output = utils.GetMergedJson(output, part)
		}
	}
	if output == nil {
		output = make(map[string]interface{}, 0)
	}
	outputJson, _ := json.Marshal(output)
	d.Set("output", string(outputJson))
	return nil
}
