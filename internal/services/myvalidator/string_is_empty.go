package myvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func StringIsEmpty() validator.String {
	return stringvalidator.LengthAtMost(0)
}
