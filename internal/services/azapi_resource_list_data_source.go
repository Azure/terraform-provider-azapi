package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceListDataSource() *schema.Resource {
	return &schema.Resource{
		Read: resourceListDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ResourceID,
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ResourceType,
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

func resourceListDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NewResourceIDSkipScopeValidation("", d.Get("parent_id").(string), d.Get("type").(string))
	if err != nil {
		return err
	}
	listUrl := strings.TrimSuffix(id.AzureResourceId, "/")

	responseBody, err := client.List(ctx, listUrl, id.ApiVersion)
	if err != nil {
		return fmt.Errorf("list resource, url: %s, error: %+v", listUrl, err)
	}

	d.SetId(listUrl)
	// #nosec G104
	d.Set("output", flattenOutput(responseBody, d.Get("response_export_values").([]interface{})))

	return nil
}
