package services

import (
	"fmt"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/services/customization"
)

// validateDataPlaneResourceWritable returns an error if the resource type has a registered
// customization but exposes no write operations (nil CreateFunc). Such types are read-only by
// design and should only be used with the data source, not the managed resource.
func validateDataPlaneResourceWritable(resourceType string) error {
	customizedResource := customization.GetCustomization(resourceType)
	if customizedResource == nil {
		return nil
	}
	if (*customizedResource).CreateFunc() == nil {
		typeBase := strings.Split(resourceType, "@")[0]
		return fmt.Errorf(`resource type %q does not support create/update/delete; use data.azapi_data_plane_resource to read it instead`, typeBase)
	}
	return nil
}
