package functions_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

type Snake2CamelFunction struct{}

func TestAccSnake2CamelFunction_basic(t *testing.T) {
	r := Snake2CamelFunction{}
	data := acceptance.BuildTestData(t, "data.azapi_functions_snake_to_camel", "test")
	data.FunctionTest(t, []resource.TestStep{
		{
			Config: r.basic(),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue(
					"output",
					knownvalue.ObjectExact(map[string]knownvalue.Check{}),
				),
			},
		},
	})
}

func (s Snake2CamelFunction) basic() string {
	return `
variable "input" {
	type = object({
		first_key = string
		second_key = object({
			third_key = string
		})
		fourth_key = list(object({
			fifth_key = string
		}))
	})
	default = {
		first_key = "value1"
		second_key = {
			third_key = "value2"
		}
		fourth_key = [
			{
				fifth_key = "value3"
			}
		]}
	})
}

output "output" {
  value = provider::azapi::snake2camel(var.input)
}
`
}
