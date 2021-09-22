package validate

import (
	"fmt"

	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/parse"
)

func ResourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := parse.ResourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
