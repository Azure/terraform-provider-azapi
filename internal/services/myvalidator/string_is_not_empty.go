package myvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func StringIsNotEmpty() validator.String {
	return stringvalidator.LengthAtLeast(1)
}
