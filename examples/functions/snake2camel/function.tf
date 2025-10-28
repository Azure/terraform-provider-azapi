locals {
  input = {
    first_thing  = string
    second_thing = string
    nested_object = object({
      nested_thing1 = string
      nested_thing2 = string
    })
  }
}

# Will output:
#
# {
#   firstThing = string
#   secondThing = string
#   nestedObject = {
#     nestedThing1 = string
#     nestedThing2 = string
#   }
# }
output "unique_name" {
  value = provider::azapi::snake2camel(local.input)
}
