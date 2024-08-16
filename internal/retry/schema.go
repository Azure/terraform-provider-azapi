package retry

import (
	"context"
	"math/big"

	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	intervalSecondsAttributeName     = "interval_seconds"
	maxIntervalSecondsAttributeName  = "max_interval_seconds"
	multiplierAttributeName          = "multiplier"
	randomizationFactorAttributeName = "randomization_factor"
	errorMessageRegexAttributeName   = "error_message_regex"
)

func SingleNestedAttribute(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "The retry block supports the following arguments:",
		Optional:            true,
		Attributes: map[string]schema.Attribute{

			intervalSecondsAttributeName: schema.Int64Attribute{
				MarkdownDescription: "The base number of seconds to wait between retries. Default is `10`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(defaultIntervalSeconds),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(120),
				},
			},

			maxIntervalSecondsAttributeName: schema.Int64Attribute{
				MarkdownDescription: "The maximum number of seconds to wait between retries. Default is `180`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(defaultMaxIntervalSeconds),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(300),
				},
			},

			multiplierAttributeName: schema.NumberAttribute{
				MarkdownDescription: "The multiplier to apply to the interval between retries. Default is `1.5`.",
				Optional:            true,
				Computed:            true,
				Default:             numberdefault.StaticBigFloat(big.NewFloat(float64(defaultMultiplier))),
			},

			randomizationFactorAttributeName: schema.NumberAttribute{
				Optional:            true,
				Computed:            true,
				Default:             numberdefault.StaticBigFloat(big.NewFloat(float64(defaultRandomizationFactor))),
				MarkdownDescription: "The randomization factor to apply to the interval between retries. The formula for the randomized interval is: `RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])`. Therefore set to zero `0.0` for no randomization. Default is `0.5`.",
			},

			errorMessageRegexAttributeName: schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A list of regular expressions to match against error messages. If any of the regular expressions match, the error is considered retryable.",
				Required:            true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsValidRegex()),
					listvalidator.UniqueValues(),
					listvalidator.SizeAtLeast(1),
				},
			},
		},
		CustomType: RetryType{
			ObjectType: types.ObjectType{
				AttrTypes: RetryValue{}.AttributeTypes(ctx),
			},
		},
	}
}
