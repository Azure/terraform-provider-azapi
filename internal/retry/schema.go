package retry

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ephemeralschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
				Required:            true,
				Description:         "A list of regular expressions to match against error messages. If any of the regular expressions match, the request will be retried.",
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
				Description:         "The base number of seconds to wait between retries.",
				MarkdownDescription: "The base number of seconds to wait between retries.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(120),
				},
				Default: int64default.StaticInt64(DefaultIntervalSeconds),
			},
			"max_interval_seconds": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The maximum number of seconds to wait between retries.",
				MarkdownDescription: "The maximum number of seconds to wait between retries.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(300),
				},
				Default: int64default.StaticInt64(DefaultMaxIntervalSeconds),
			},
			"multiplier": schema.Float64Attribute{
				DeprecationMessage:  "The `multiplier` attribute is deprecated and will be removed in a future version.",
				Optional:            true,
				Computed:            true,
				Description:         "The multiplier to apply to the interval between retries.",
				MarkdownDescription: "The multiplier to apply to the interval between retries.",
				Default:             float64default.StaticFloat64(DefaultMultiplier),
			},
			"randomization_factor": schema.Float64Attribute{
				DeprecationMessage:  "The `randomization_factor` attribute is deprecated and will be removed in a future version.",
				Optional:            true,
				Computed:            true,
				Description:         "The randomization factor to apply to the interval between retries. The formula for the randomized interval is: `RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])`. Therefore set to zero `0.0` for no randomization.",
				MarkdownDescription: "The randomization factor to apply to the interval between retries. The formula for the randomized interval is: `RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])`. Therefore set to zero `0.0` for no randomization.",
				Default:             float64default.StaticFloat64(DefaultRandomizationFactor),
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

// Same as Retry() but for data source.
func RetryDsSchema(ctx context.Context) dsschema.Attribute {
	return dsschema.SingleNestedAttribute{
		Attributes: map[string]dsschema.Attribute{
			"error_message_regex": dsschema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				Description:         "A list of regular expressions to match against error messages. If any of the regular expressions match, the request will be retried.",
				MarkdownDescription: "A list of regular expressions to match against error messages. If any of the regular expressions match, the request will be retried.",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsValidRegex()),
					listvalidator.UniqueValues(),
					listvalidator.SizeAtLeast(1),
				},
			},
			"interval_seconds": dsschema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The base number of seconds to wait between retries.",
				MarkdownDescription: "The base number of seconds to wait between retries.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(120),
				},
			},
			"max_interval_seconds": dsschema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The maximum number of seconds to wait between retries.",
				MarkdownDescription: "The maximum number of seconds to wait between retries.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(300),
				},
			},
			"multiplier": dsschema.Float64Attribute{
				DeprecationMessage:  "The `multiplier` attribute is deprecated and will be removed in a future version.",
				Optional:            true,
				Computed:            true,
				Description:         "The multiplier to apply to the interval between retries.",
				MarkdownDescription: "The multiplier to apply to the interval between retries.",
			},
			"randomization_factor": dsschema.Float64Attribute{
				DeprecationMessage:  "The `randomization_factor` attribute is deprecated and will be removed in a future version.",
				Optional:            true,
				Computed:            true,
				Description:         "The randomization factor to apply to the interval between retries. The formula for the randomized interval is: `RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])`. Therefore set to zero `0.0` for no randomization.",
				MarkdownDescription: "The randomization factor to apply to the interval between retries. The formula for the randomized interval is: `RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])`. Therefore set to zero `0.0` for no randomization.",
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

// Same as Retry() but for ephemeral resource.
func RetryEphemeralSchema(ctx context.Context) ephemeralschema.Attribute {
	return ephemeralschema.SingleNestedAttribute{
		Attributes: map[string]ephemeralschema.Attribute{
			"error_message_regex": ephemeralschema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				Description:         "A list of regular expressions to match against error messages. If any of the regular expressions match, the request will be retried.",
				MarkdownDescription: "A list of regular expressions to match against error messages. If any of the regular expressions match, the request will be retried.",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsValidRegex()),
					listvalidator.UniqueValues(),
					listvalidator.SizeAtLeast(1),
				},
			},
			"interval_seconds": ephemeralschema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The base number of seconds to wait between retries.",
				MarkdownDescription: "The base number of seconds to wait between retries.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(120),
				},
			},
			"max_interval_seconds": ephemeralschema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The maximum number of seconds to wait between retries.",
				MarkdownDescription: "The maximum number of seconds to wait between retries.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(300),
				},
			},
			"multiplier": ephemeralschema.Float64Attribute{
				DeprecationMessage:  "The `multiplier` attribute is deprecated and will be removed in a future version.",
				Optional:            true,
				Computed:            true,
				Description:         "The multiplier to apply to the interval between retries.",
				MarkdownDescription: "The multiplier to apply to the interval between retries.",
			},
			"randomization_factor": ephemeralschema.Float64Attribute{
				DeprecationMessage:  "The `randomization_factor` attribute is deprecated and will be removed in a future version.",
				Optional:            true,
				Computed:            true,
				Description:         "The randomization factor to apply to the interval between retries. The formula for the randomized interval is: `RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])`. Therefore set to zero `0.0` for no randomization.",
				MarkdownDescription: "The randomization factor to apply to the interval between retries. The formula for the randomized interval is: `RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])`. Therefore set to zero `0.0` for no randomization.",
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
