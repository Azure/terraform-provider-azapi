package tf

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IDValidationFunc func(id string) error

func DefaultImporter(validateFunc IDValidationFunc) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			log.Printf("[DEBUG] Importing Resource - parsing %q", d.Id())

			if err := validateFunc(d.Id()); err != nil {
				return []*schema.ResourceData{d}, fmt.Errorf("parsing Resource ID %q: %+v", d.Id(), err)
			}

			return []*schema.ResourceData{d}, nil
		},
	}
}

func ImportAsExistsError(resourceName, id string) error {
	msg := "A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information."
	return fmt.Errorf(msg, id, resourceName)
}
