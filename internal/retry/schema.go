package retry

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func RetrySchema(ctx context.Context) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"error_message_regex": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A list of regular expressions to match against error messages. If any of the regular expressions match, the request will be retried.",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsValidRegex()),
					listvalidator.UniqueValues(),
					listvalidator.SizeAtLeast(1),
				},
			},
			"interval_seconds": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The base number of seconds to wait between retries. Default is `10`.",
				MarkdownDescription: "The base number of seconds to wait between retries. Default is `10`.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(120),
				},
				Default: int64default.StaticInt64(10),
			},
			"max_interval_seconds": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The maximum number of seconds to wait between retries. Default is `180`.",
				MarkdownDescription: "The maximum number of seconds to wait between retries. Default is `180`.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(300),
				},
				Default: int64default.StaticInt64(180),
			},
			"multiplier": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The multiplier to apply to the interval between retries. Default is `1.5`.",
				MarkdownDescription: "The multiplier to apply to the interval between retries. Default is `1.5`.",
				Default:             float64default.StaticFloat64(1.5),
			},
			"randomization_factor": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The randomization factor to apply to the interval between retries. The formula for the randomized interval is: `RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])`. Therefore set to zero `0.0` for no randomization. Default is `0.5`.",
				MarkdownDescription: "The randomization factor to apply to the interval between retries. The formula for the randomized interval is: `RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])`. Therefore set to zero `0.0` for no randomization. Default is `0.5`.",
				Default:             float64default.StaticFloat64(0.5),
			},
			"response_is_nil": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If set to `true`, the retry will be triggered when the response is `nil`. Default is `true`.",
				MarkdownDescription: "If set to `true`, the retry will be triggered when the response is `nil`. Default is `true`.",
				Default:             booldefault.StaticBool(true),
			},
			"status_forbidden": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If set to `true`, the retry will be triggered when the status code is `403`. Default is `false`.",
				MarkdownDescription: "If set to `true`, the retry will be triggered when the status code is `403`. Default is `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"status_not_found": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If set to `true`, the retry will be triggered when the status code is `404`. Default is `false`.",
				MarkdownDescription: "If set to `true`, the retry will be triggered when the status code is `404`. Default is `false`.",
				Default:             booldefault.StaticBool(false),
			},
		},
		CustomType: RetryType{
			ObjectType: types.ObjectType{
				AttrTypes: RetryValue{}.AttributeTypes(ctx),
			},
		},
		Optional:            true,
		Description:         "The retry object supports the following attributes:",
		MarkdownDescription: "The retry object supports the following attributes:",
	}
}
