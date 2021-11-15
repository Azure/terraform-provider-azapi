package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
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

func AzureResourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	r := regexp.MustCompile("^http[s]?:.*")
	if r.MatchString(v) {
		errors = append(errors, fmt.Errorf("expected %q not to contain protocol", key))
	}

	if _, err := resourceids.ParseAzureResourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

func ResourceType(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if v == "" {
		return nil, []error{fmt.Errorf("expected %q to not be an empty string, got %v", k, i)}
	}

	parts := strings.Split(v, "@")
	if len(parts) != 2 {
		return nil, []error{fmt.Errorf("expected %q to be <resource-type>@<api-version>", k)}
	}

	return nil, nil
}
