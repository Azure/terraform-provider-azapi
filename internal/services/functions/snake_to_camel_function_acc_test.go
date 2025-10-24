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
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"firstKey": knownvalue.StringExact("value1"),
					}),
				),
			},
		},
	})
}

func TestAccSnake2CamelFunction_partialUnkown(t *testing.T) {
	r := Snake2CamelFunction{}
	data := acceptance.BuildTestData(t, "data.azapi_functions_snake_to_camel", "test")
	data.FunctionTest(t, []resource.TestStep{
		{
			Config: r.partialUnknown(),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue(
					"output",
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"firstKey": knownvalue.StringExact("unknown"),
					}),
				),
			},
		},
	})
}
func TestAccSnake2CamelFunction_fullyUnkown(t *testing.T) {
	r := Snake2CamelFunction{}
	data := acceptance.BuildTestData(t, "data.azapi_functions_snake_to_camel", "test")
	data.FunctionTest(t, []resource.TestStep{
		{
			Config: r.fullyUnknown(),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue(
					"output",
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"firstKey": knownvalue.StringExact("unknown"),
					}),
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
	})
	default = {
		first_key = "value1"
	}
}

output "output" {
  value = provider::azapi::snake2camel(var.input)
}
`
}

func (s Snake2CamelFunction) partialUnknown() string {
	return `
resource "terraform_data" "test" {
	input = "unknown"
}

output "output" {
  value	 = provider::azapi::snake2camel({
	  first_key = terraform_data.test.output
	})
}
`
}

func (s Snake2CamelFunction) fullyUnknown() string {
	return `
resource "terraform_data" "test" {
	input = {
		first_key = "unknown"
	}
}

output "output" {
  value	 = provider::azapi::snake2camel(terraform_data.test.output)
}
`
}

func TestAccSnake2CamelFunction_nestedObjectWithKnownValues(t *testing.T) {
	r := Snake2CamelFunction{}
	data := acceptance.BuildTestData(t, "data.azapi_functions_snake_to_camel", "test")
	data.FunctionTest(t, []resource.TestStep{
		{
			Config: r.nestedObjectWithKnownValues(),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue(
					"output",
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"firstKey":  knownvalue.StringExact("value1"),
						"secondKey": knownvalue.StringExact("value2"),
						"nestedObject": knownvalue.ObjectExact(map[string]knownvalue.Check{
							"childKey":     knownvalue.StringExact("childValue"),
							"anotherChild": knownvalue.StringExact("anotherValue"),
						}),
					}),
				),
			},
		},
	})
}

func TestAccSnake2CamelFunction_nestedObjectWithUnknownChild(t *testing.T) {
	r := Snake2CamelFunction{}
	data := acceptance.BuildTestData(t, "data.azapi_functions_snake_to_camel", "test")
	data.FunctionTest(t, []resource.TestStep{
		{
			Config: r.nestedObjectWithUnknownChild(),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue(
					"output",
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"firstKey": knownvalue.StringExact("knownValue"),
						"nestedObject": knownvalue.ObjectExact(map[string]knownvalue.Check{
							"childKey":     knownvalue.StringExact("unknown"),
							"anotherChild": knownvalue.StringExact("knownChild"),
						}),
					}),
				),
			},
		},
	})
}

func TestAccSnake2CamelFunction_deeplyNestedObject(t *testing.T) {
	r := Snake2CamelFunction{}
	data := acceptance.BuildTestData(t, "data.azapi_functions_snake_to_camel", "test")
	data.FunctionTest(t, []resource.TestStep{
		{
			Config: r.deeplyNestedObject(),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue(
					"output",
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"topLevel": knownvalue.ObjectExact(map[string]knownvalue.Check{
							"midLevel": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"deepLevel": knownvalue.StringExact("deepValue"),
							}),
						}),
					}),
				),
			},
		},
	})
}

func TestAccSnake2CamelFunction_nestedObjectWithList(t *testing.T) {
	r := Snake2CamelFunction{}
	data := acceptance.BuildTestData(t, "data.azapi_functions_snake_to_camel", "test")
	data.FunctionTest(t, []resource.TestStep{
		{
			Config: r.nestedObjectWithList(),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue(
					"output",
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"myList": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"itemKey": knownvalue.StringExact("item1"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"itemKey": knownvalue.StringExact("item2"),
							}),
						}),
					}),
				),
			},
		},
	})
}

func TestAccSnake2CamelFunction_mixedKnownUnknownNested(t *testing.T) {
	r := Snake2CamelFunction{}
	data := acceptance.BuildTestData(t, "data.azapi_functions_snake_to_camel", "test")
	data.FunctionTest(t, []resource.TestStep{
		{
			Config: r.mixedKnownUnknownNested(),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue(
					"output",
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"knownKey":   knownvalue.StringExact("knownValue"),
						"unknownKey": knownvalue.StringExact("unknown"),
						"nestedObject": knownvalue.ObjectExact(map[string]knownvalue.Check{
							"knownChild":   knownvalue.StringExact("knownChildValue"),
							"unknownChild": knownvalue.StringExact("unknownChild"),
						}),
					}),
				),
			},
		},
	})
}

func (s Snake2CamelFunction) nestedObjectWithKnownValues() string {
	return `
variable "input" {
	type = object({
		first_key = string
		second_key = string
		nested_object = object({
			child_key = string
			another_child = string
		})
	})
	default = {
		first_key = "value1"
		second_key = "value2"
		nested_object = {
			child_key = "childValue"
			another_child = "anotherValue"
		}
	}
}

output "output" {
  value = provider::azapi::snake2camel(var.input)
}
`
}

func (s Snake2CamelFunction) nestedObjectWithUnknownChild() string {
	return `
resource "terraform_data" "test" {
	input = "unknown"
}

output "output" {
  value = provider::azapi::snake2camel({
	  first_key = "knownValue"
	  nested_object = {
		  child_key = terraform_data.test.output
		  another_child = "knownChild"
	  }
  })
}
`
}

func (s Snake2CamelFunction) deeplyNestedObject() string {
	return `
variable "input" {
	type = object({
		top_level = object({
			mid_level = object({
				deep_level = string
			})
		})
	})
	default = {
		top_level = {
			mid_level = {
				deep_level = "deepValue"
			}
		}
	}
}

output "output" {
  value = provider::azapi::snake2camel(var.input)
}
`
}

func (s Snake2CamelFunction) nestedObjectWithList() string {
	return `
variable "input" {
	type = object({
		my_list = list(object({
			item_key = string
		}))
	})
	default = {
		my_list = [
			{
				item_key = "item1"
			},
			{
				item_key = "item2"
			}
		]
	}
}

output "output" {
  value = provider::azapi::snake2camel(var.input)
}
`
}

func (s Snake2CamelFunction) mixedKnownUnknownNested() string {
	return `
resource "terraform_data" "test" {
	input = "unknown"
}

resource "terraform_data" "test2" {
	input = "unknownChild"
}

output "output" {
  value = provider::azapi::snake2camel({
	  known_key = "knownValue"
	  unknown_key = terraform_data.test.output
	  nested_object = {
		  known_child = "knownChildValue"
		  unknown_child = terraform_data.test2.output
	  }
  })
}
`
}
