package acceptance

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// RequiresImportErrorStep returns a Test Step which expects a Requires Import
// error to be returned when running this step
func (td TestData) RequiresImportErrorStep(configBuilder func(data TestData) string) resource.TestStep {
	config := configBuilder(td)
	return resource.TestStep{
		Config:      config,
		ExpectError: RequiresImportError(td.ResourceType),
	}
}

func RequiresImportError(_ string) *regexp.Regexp {
	message := "Resource already exists"
	return regexp.MustCompile(message)
}

// ImportStep returns a Test Step which Imports the Resource, optionally
// ignoring any fields which may not be imported (for example, as they're
// not returned from the API)
func (td TestData) ImportStep(ignore ...string) resource.TestStep {
	return td.ImportStepFor(td.ResourceName, ignore...)
}

// ImportStep returns a Test Step which Imports the Resource, optionally
// ignoring any fields which may not be imported (for example, as they're
// not returned from the API)
func (td TestData) ImportStepWithImportStateIdFunc(importStateIdFunc resource.ImportStateIdFunc, ignore ...string) resource.TestStep {
	resourceName := td.ResourceName
	step := resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: importStateIdFunc,
	}

	if len(ignore) > 0 {
		step.ImportStateVerifyIgnore = ignore
	}

	return step
}

// ImportStepFor returns a Test Step which Imports a given resource by name,
// optionally ignoring any fields which may not be imported (for example, as they're
// not returned from the API)
func (td TestData) ImportStepFor(resourceName string, ignore ...string) resource.TestStep {
	if strings.HasPrefix(resourceName, "data.") {
		return resource.TestStep{
			ResourceName: resourceName,
			SkipFunc: func() (bool, error) {
				return false, fmt.Errorf("Data Sources (%q) do not support import - remove the ImportStep / ImportStepFor`", resourceName)
			},
		}
	}

	step := resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateVerify: true,
	}

	if len(ignore) > 0 {
		step.ImportStateVerifyIgnore = ignore
	}

	return step
}
