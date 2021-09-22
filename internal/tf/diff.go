package tf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

func SuppressJsonOrderingDifference(_, old, new string, _ *schema.ResourceData) bool {
	return utils.NormalizeJson(old) == utils.NormalizeJson(new)
}
