package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

func ResourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if v == "/" {
		return
	}

	r := regexp.MustCompile("^http[s]?:.*")
	if r.MatchString(v) {
		errors = append(errors, fmt.Errorf("expected %q not to contain protocol, got %q", key, v))
	}
	r = regexp.MustCompile(".*api-version=.*")
	if r.MatchString(v) {
		errors = append(errors, fmt.Errorf("expected %q not to contain api-version query parameter, got %q", key, v))
	}

	if strings.Contains(v, "|") {
		errors = append(errors, fmt.Errorf("expected %q not to contain character '|', got %q. This usually indicates a synthetic resource ID created by the AzureRM provider, which is not the expected format for an Azure Resource ID. Please use the actual resource ID instead", key, v))
	}

	if _, err := arm.ParseResourceID(v); err != nil {
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
