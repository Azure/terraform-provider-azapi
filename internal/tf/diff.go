package tf

import (
	"github.com/Azure/terraform-provider-azurerm-restapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SuppressJsonOrderingDifference(_, old, new string, _ *schema.ResourceData) bool {
	return utils.NormalizeJson(old) == utils.NormalizeJson(new)
}
