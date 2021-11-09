package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/azure"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/azure/identity"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/azure/location"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/azure/tags"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/clients"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/parse"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/validate"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/tf"
	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

func ResourceAzureGenericResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureGenericResourceCreateUpdate,
		Read:   resourceAzureGenericResourceRead,
		Update: resourceAzureGenericResourceCreateUpdate,
		Delete: resourceAzureGenericResourceDelete,

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
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AzureResourceID,
			},

			"api_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": location.SchemaLocation(),

			"identity": identity.SchemaIdentity(),

			"body": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: tf.SuppressJsonOrderingDifference,
			},

			"create_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "PUT",
				ValidateFunc: validation.StringInSlice([]string{
					http.MethodPost,
					http.MethodPut,
				}, false),
			},

			"update_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "PUT",
				ValidateFunc: validation.StringInSlice([]string{
					http.MethodPost,
					http.MethodPut,
					// http.MethodPatch, patch is not supported yet
				}, false),
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

			"tags": tags.SchemaTags(),
		},

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			if d.HasChange("identity") || d.HasChange("tags") || d.HasChange("response_export_values") {
				d.SetNewComputed("output")
			}
			old, new := d.GetChange("body")
			if utils.NormalizeJson(old) != utils.NormalizeJson(new) {
				d.SetNewComputed("output")
			}

			var body interface{}
			err := json.Unmarshal([]byte(d.Get("body").(string)), &body)
			if err != nil {
				return err
			}
			props := []string{"identity", "location", "tags"}
			for _, prop := range props {
				if _, ok := d.GetOk(prop); ok {
					if bodyMap, ok := body.(map[string]interface{}); ok {
						if bodyMap[prop] != nil {
							return fmt.Errorf("can't specify both property `%[1]s` and `%[1]s` in `body`", prop)
						}
					}
				}
			}

			if meta.(*clients.Client).Features.SchemaValidationEnabled {
				if value, ok := d.GetOk("tags"); ok {
					bodyWithTags := tags.ExpandTags(value.(map[string]interface{}))
					body = utils.GetMergedJson(body, bodyWithTags)
				}
				if value, ok := d.GetOk("location"); ok {
					bodyWithLocation := location.ExpandLocation(value.(string))
					body = utils.GetMergedJson(body, bodyWithLocation)
				}
				if value, ok := d.GetOk("identity"); ok {
					bodyWithIdentity, err := identity.ExpandIdentity(value.([]interface{}))
					if err != nil {
						return err
					}
					body = utils.GetMergedJson(body, bodyWithIdentity)
				}
				if err := schemaValidation(d.Get("resource_id").(string), d.Get("api_version").(string), body); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func resourceAzureGenericResourceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewResourceID(d.Get("resource_id").(string), d.Get("api_version").(string))

	if d.IsNewResource() {
		existing, response, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
		if err != nil {
			if response.StatusCode != http.StatusNotFound {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if len(utils.GetId(existing)) > 0 {
			return tf.ImportAsExistsError("azurermg_resource", id.ID())
		}
	}

	var requestBody interface{}
	err := json.Unmarshal([]byte(d.Get("body").(string)), &requestBody)
	if err != nil {
		return err
	}

	props := []string{"identity", "location", "tags"}
	for _, prop := range props {
		if _, ok := d.GetOk(prop); ok {
			if bodyMap, ok := requestBody.(map[string]interface{}); ok {
				if bodyMap[prop] != nil {
					return fmt.Errorf("can't specify both property `%[1]s` and `%[1]s` in `body`", prop)
				}
			}
		}
	}

	if value, ok := d.GetOk("tags"); ok {
		bodyWithTags := tags.ExpandTags(value.(map[string]interface{}))
		requestBody = utils.GetMergedJson(requestBody, bodyWithTags)
	}
	if value, ok := d.GetOk("location"); ok {
		bodyWithLocation := location.ExpandLocation(value.(string))
		requestBody = utils.GetMergedJson(requestBody, bodyWithLocation)
	}
	if value, ok := d.GetOk("identity"); ok {
		bodyWithIdentity, err := identity.ExpandIdentity(value.([]interface{}))
		if err != nil {
			return err
		}
		requestBody = utils.GetMergedJson(requestBody, bodyWithIdentity)
	}

	if meta.(*clients.Client).Features.SchemaValidationEnabled {
		if err := schemaValidation(id.AzureResourceId, id.ApiVersion, requestBody); err != nil {
			return err
		}
	}

	var method string
	switch {
	case d.IsNewResource():
		method = d.Get("create_method").(string)
	default:
		method = d.Get("update_method").(string)
	}

	j, _ := json.Marshal(requestBody)
	log.Printf("[INFO] request body: %v\n", string(j))
	_, _, err = client.CreateUpdate(ctx, id.AzureResourceId, id.ApiVersion, requestBody, method)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAzureGenericResourceRead(d, meta)
}

func resourceAzureGenericResourceRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[INFO] Error reading %q - removing from state", id.ID())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	d.Set("resource_id", id.AzureResourceId)
	d.Set("api_version", id.ApiVersion)

	bodyJson := d.Get("body").(string)
	var requestBody interface{}
	err = json.Unmarshal([]byte(bodyJson), &requestBody)
	if err != nil {
		if len(bodyJson) == 0 {
			// handle import case, body defaults to response body
			requestBody = utils.GetIgnoredJson(responseBody, getUnsupportedProperties())
		} else {
			return err
		}
	}
	data, err := json.Marshal(utils.GetUpdatedJson(requestBody, responseBody))
	if err != nil {
		return err
	}
	d.Set("body", string(data))
	d.Set("tags", tags.FlattenTags(responseBody))
	d.Set("location", location.FlattenLocation(responseBody))
	d.Set("identity", identity.FlattenIdentity(responseBody))

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

func resourceAzureGenericResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceID(d.Id())
	if err != nil {
		return err
	}

	_, _, err = client.Delete(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}

func schemaValidation(resourceId, apiVersion string, body interface{}) error {
	if len(resourceId) == 0 {
		log.Println("[INFO] `resource_id` will be known after apply, skip schema validation on plan phase.")
		return nil
	}

	resourceType := utils.GetResourceType(resourceId)
	log.Printf("[INFO] prepare validation for resource type: %s", resourceType)
	versions := azure.GetApiVersions(resourceType)
	isVersionValid := false
	for _, version := range versions {
		if version == apiVersion {
			isVersionValid = true
			break
		}
	}
	if !isVersionValid {
		return fmt.Errorf("the `api_version` is invalid. The supported versions are [%s]\n", strings.Join(versions, ", "))
	}

	resourceDef, err := azure.GetResourceDefinition(resourceType, apiVersion)
	if err == nil && resourceDef != nil {
		errors := (*resourceDef).Validate(utils.NormalizeObject(body), "")
		if len(errors) != 0 {
			errorMsg := "the `body` is invalid: \n"
			for _, err := range errors {
				errorMsg += fmt.Sprintf("%s\n", err.Error())
			}
			return fmt.Errorf(errorMsg)
		}
	} else {
		log.Printf("[ERROR] load embedded schema: %+v\n", err)
	}
	return nil
}
