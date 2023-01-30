package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/resourceName"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/locks"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/internal/services/validate"
	"github.com/Azure/terraform-provider-azapi/internal/tf"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAzApiResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzApiResourceCreateUpdate,
		Read:   resourceAzApiResourceRead,
		Update: resourceAzApiResourceCreateUpdate,
		Delete: resourceAzApiResourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				log.Printf("[DEBUG] Importing Resource - parsing %q", d.Id())

				input := d.Id()
				idUrl, err := url.Parse(input)
				if err != nil {
					return []*schema.ResourceData{d}, fmt.Errorf("parsing Resource ID %q: %+v", input, err)
				}
				apiVersion := idUrl.Query().Get("api-version")
				if len(apiVersion) == 0 {
					resourceType := utils.GetResourceType(input)
					apiVersions := azure.GetApiVersions(resourceType)
					if len(apiVersions) != 0 {
						input = fmt.Sprintf("%s?api-version=%s", input, apiVersions[len(apiVersions)-1])
					}
				}

				id, err := parse.ResourceIDWithApiVersion(input)
				if err != nil {
					return []*schema.ResourceData{d}, fmt.Errorf("parsing Resource ID %q: %+v", d.Id(), err)
				}
				// override the id to remove the api-version
				d.SetId(id.ID())
				// #nosec G104
				d.Set("type", fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion))
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": resourceName.SchemaResourceNameOC(),

			"removing_special_chars": resourceName.SchemaResourceNameRemovingSpecialCharacters(),

			"parent_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ResourceID,
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ResourceType,
			},

			"location": location.SchemaLocationOC(),

			"identity": identity.SchemaIdentity(),

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

			"schema_validation_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"output": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaTagsOC(),
		},

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			if d.HasChange("identity") || d.HasChange("tags") || d.HasChange("response_export_values") {
				// #nosec G104
				d.SetNewComputed("output")
			}
			old, new := d.GetChange("body")
			if utils.NormalizeJson(old) != utils.NormalizeJson(new) {
				// #nosec G104
				d.SetNewComputed("output")
			}

			parentId := d.Get("parent_id").(string)
			resourceType := d.Get("type").(string)
			var assignedName string

			config := d.GetRawConfig()
			if !isConfigExist(config, "name") && len(meta.(*clients.Client).Features.DefaultNaming) == 0 {
				return fmt.Errorf("resource name can't be empty, either specifying a default name or a resource name")
			}

			// body refers other resource, can't be verified during plan
			if len(d.Get("body").(string)) == 0 {
				return nil
			}

			var body map[string]interface{}
			err := json.Unmarshal([]byte(d.Get("body").(string)), &body)
			if err != nil {
				return err
			}

			props := []string{"identity", "location", "tags"}
			for _, prop := range props {
				if isConfigExist(config, prop) && body[prop] != nil {
					return fmt.Errorf("can't specify both property `%[1]s` and `%[1]s` in `body`", prop)
				}
			}

			azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(d.Get("type").(string))
			if err != nil {
				return err
			}
			resourceDef, _ := azure.GetResourceDefinition(azureResourceType, apiVersion)
			if !isConfigExist(config, "tags") && body["tags"] == nil && len(meta.(*clients.Client).Features.DefaultTags) != 0 {
				if isResourceHasProperty(resourceDef, "tags") {
					body["tags"] = meta.(*clients.Client).Features.DefaultTags
					currentTags := d.Get("tags")
					defaultTags := meta.(*clients.Client).Features.DefaultTags
					if !reflect.DeepEqual(currentTags, defaultTags) {
						// #nosec G104
						d.SetNew("tags", defaultTags)
					}
				}
			}

			if !isConfigExist(config, "name") && len(meta.(*clients.Client).Features.DefaultNaming) != 0 {
				currentName := d.Get("name").(string)
				defaultName := meta.(*clients.Client).Features.DefaultNaming
				assignedName = defaultName
				if currentName != defaultName {
					// #nosec G104
					d.SetNew("name", assignedName)
				}
			}

			if value, ok := d.GetOk("name"); ok && isConfigExist(config, "name") {
				currentName := d.Get("name").(string)
				assignedName = value.(string)

				if len(meta.(*clients.Client).Features.DefaultNamingPrefix) != 0 {
					assignedName = meta.(*clients.Client).Features.DefaultNamingPrefix + assignedName
				}
				if len(meta.(*clients.Client).Features.DefaultNamingSuffix) != 0 {
					assignedName += meta.(*clients.Client).Features.DefaultNamingSuffix
				}
				if _, ok := d.GetOk("removing_special_chars"); ok && isConfigExist(config, "removing_special_chars") {
					assignedName = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(assignedName, "")
				}
				if currentName != assignedName {
					// #nosec G104
					d.SetNew("name", assignedName)
				}
			}

			// verify parent_id when it's known
			if len(parentId) > 0 {
				_, err := parse.NewResourceID(assignedName, parentId, resourceType)
				if err != nil {
					return err
				}
			}

			if !isConfigExist(config, "location") && body["location"] == nil && len(meta.(*clients.Client).Features.DefaultLocation) != 0 {
				if isResourceHasProperty(resourceDef, "location") {
					body["location"] = meta.(*clients.Client).Features.DefaultLocation
					currentLocation := d.Get("location").(string)
					defaultLocation := meta.(*clients.Client).Features.DefaultLocation
					if location.Normalize(currentLocation) != location.Normalize(defaultLocation) {
						// #nosec G104
						d.SetNew("location", defaultLocation)
					}
				}
			}

			if d.Get("schema_validation_enabled").(bool) {
				if value, ok := d.GetOk("tags"); ok && isConfigExist(config, "tags") {
					tagsModel := tags.ExpandTags(value.(map[string]interface{}))
					if len(tagsModel) != 0 {
						body["tags"] = tagsModel
					}
				}
				if value, ok := d.GetOk("location"); ok && isConfigExist(config, "location") {
					body["location"] = location.Normalize(value.(string))
				}
				if value, ok := d.GetOk("identity"); ok && isConfigExist(config, "identity") {
					identityModel, err := identity.ExpandIdentity(value.([]interface{}))
					if err != nil {
						return err
					}
					if identityModel != nil {
						body["identity"] = identityModel
					}
				}
				if err := schemaValidation(azureResourceType, apiVersion, resourceDef, body); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func resourceAzApiResourceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	config := d.GetRawConfig()
	var resourceName string

	if !isConfigExist(config, "name") && len(meta.(*clients.Client).Features.DefaultNaming) == 0 {
		return fmt.Errorf("either a default name or a user assigned name in config file needs to be assigned")
	}

	if value, ok := d.GetOk("name"); ok && isConfigExist(config, "name") {
		resourceName = value.(string)
	} else {
		resourceName = meta.(*clients.Client).Features.DefaultNaming
	}

	id, err := parse.NewResourceID(resourceName, d.Get("parent_id").(string), d.Get("type").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		_, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
		if err == nil {
			return tf.ImportAsExistsError("azapi_resource", id.ID())
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

	props := []string{"identity", "location", "tags"}
	for _, prop := range props {
		if isConfigExist(config, prop) && body[prop] != nil {
			return fmt.Errorf("can't specify both property `%[1]s` and `%[1]s` in `body`", prop)
		}
	}

	if !isConfigExist(config, "tags") && body["tags"] == nil && len(meta.(*clients.Client).Features.DefaultTags) != 0 {
		if isResourceHasProperty(id.ResourceDef, "tags") {
			body["tags"] = meta.(*clients.Client).Features.DefaultTags
		}
	}

	if !isConfigExist(config, "location") && body["location"] == nil && len(meta.(*clients.Client).Features.DefaultLocation) != 0 {
		if isResourceHasProperty(id.ResourceDef, "location") {
			body["location"] = meta.(*clients.Client).Features.DefaultLocation
		}
	}

	if value, ok := d.GetOk("tags"); ok && isConfigExist(config, "tags") {
		tagsModel := tags.ExpandTags(value.(map[string]interface{}))
		if len(tagsModel) != 0 {
			body["tags"] = tagsModel
		}
	}
	if value, ok := d.GetOk("location"); ok && isConfigExist(config, "location") {
		body["location"] = location.Normalize(value.(string))
	}
	if value, ok := d.GetOk("identity"); ok && isConfigExist(config, "identity") {
		identityModel, err := identity.ExpandIdentity(value.([]interface{}))
		if err != nil {
			return err
		}
		if identityModel != nil {
			body["identity"] = identityModel
		}
	}

	if d.Get("schema_validation_enabled").(bool) {
		if err := schemaValidation(id.AzureResourceType, id.ApiVersion, id.ResourceDef, body); err != nil {
			return err
		}
	}

	j, _ := json.Marshal(body)
	log.Printf("[INFO] request body: %v\n", string(j))

	for _, id := range d.Get("locks").([]interface{}) {
		locks.ByID(id.(string))
		defer locks.UnlockByID(id.(string))
	}

	_, err = client.CreateOrUpdate(ctx, id.AzureResourceId, id.ApiVersion, body)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAzApiResourceRead(d, meta)
}

func resourceAzApiResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceIDWithResourceType(d.Id(), d.Get("type").(string))
	if err != nil {
		return err
	}

	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
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

	// if it's imported
	if len(d.Get("name").(string)) == 0 {
		if id.ResourceDef != nil {
			writeOnlyBody := (*id.ResourceDef).GetWriteOnly(responseBody)
			if bodyMap, ok := writeOnlyBody.(map[string]interface{}); ok {
				delete(bodyMap, "location")
				delete(bodyMap, "tags")
				delete(bodyMap, "name")
				delete(bodyMap, "identity")
				writeOnlyBody = bodyMap
			}
			data, err := json.Marshal(writeOnlyBody)
			if err != nil {
				return err
			}
			// #nosec G104
			d.Set("body", string(data))
		}
		// #nosec G104
		d.Set("ignore_casing", false)
		// #nosec G104
		d.Set("ignore_missing_property", true)
		// #nosec G104
		d.Set("schema_validation_enabled", true)
		// #nosec G104
		d.Set("removing_special_chars", false)
	} else {
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
	}

	// #nosec G104
	d.Set("name", id.Name)
	// #nosec G104
	d.Set("parent_id", id.ParentId)
	// #nosec G104
	d.Set("type", fmt.Sprintf("%s@%s", id.AzureResourceType, id.ApiVersion))

	if bodyMap, ok := responseBody.(map[string]interface{}); ok {
		// #nosec G104
		d.Set("tags", tags.FlattenTags(bodyMap["tags"]))
		// #nosec G104
		d.Set("location", bodyMap["location"])
		// #nosec G104
		d.Set("identity", identity.FlattenIdentity(bodyMap["identity"]))
	}

	// #nosec G104
	d.Set("output", flattenOutput(responseBody, d.Get("response_export_values").([]interface{})))
	return nil
}

func resourceAzApiResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := tf.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceIDWithResourceType(d.Id(), d.Get("type").(string))
	if err != nil {
		return err
	}

	for _, id := range d.Get("locks").([]interface{}) {
		locks.ByID(id.(string))
		defer locks.UnlockByID(id.(string))
	}

	_, err = client.Delete(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			return nil
		}
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}

func isConfigExist(config cty.Value, path string) bool {
	if config.CanIterateElements() {
		configMap := config.AsValueMap()
		if value, ok := configMap[path]; ok {
			if value.Type().IsListType() {
				return len(value.AsValueSlice()) != 0
			}
			return !value.IsNull()
		}
	}
	return false
}
