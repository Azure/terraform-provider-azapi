package retry

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.ObjectTypable  = RetryType{}
	_ basetypes.ObjectValuable = RetryValue{}
)

const (
	defaultIntervalSeconds     = 10
	defaultMaxIntervalSeconds  = 180
	defaultMultiplier          = 1.5
	defaultRandomizationFactor = 0.5
)

type RetryType struct {
	basetypes.ObjectType
}

func (t RetryType) ValueType(ctx context.Context) attr.Value {
	return RetryValue{}
}

func (t RetryType) String() string {
	return "retry.RetryableErrorsType"
}

func (t RetryType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	errorMessageRegexAttribute, ok := attributes[errorMessageRegexAttributeName]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			fmt.Sprintf(`%s is missing from object`, errorMessageRegexAttributeName))

		return nil, diags
	}

	errorMessageRegexVal, ok := errorMessageRegexAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`%s expected to be basetypes.ListValue, was: %T`, errorMessageRegexAttributeName, errorMessageRegexAttribute))
	}

	intervalSecondsAttribute, ok := attributes[intervalSecondsAttributeName]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			fmt.Sprintf(`%s is missing from object`, intervalSecondsAttributeName))

		return nil, diags
	}

	intervalSecondsVal, ok := intervalSecondsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`%s expected to be basetypes.Int64Value, was: %T`, intervalSecondsAttributeName, intervalSecondsAttribute))
	}

	maxIntervalSecondsAttribute, ok := attributes[maxIntervalSecondsAttributeName]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			fmt.Sprintf(`%s is missing from object`, maxIntervalSecondsAttributeName))

		return nil, diags
	}

	maxIntervalSecondsVal, ok := maxIntervalSecondsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`%s expected to be basetypes.Int64Value, was: %T`, maxIntervalSecondsAttributeName, maxIntervalSecondsAttribute))
	}

	multiplierAttribute, ok := attributes[multiplierAttributeName]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			fmt.Sprintf(`%s is missing from object`, multiplierAttributeName))

		return nil, diags
	}

	multiplierVal, ok := multiplierAttribute.(basetypes.NumberValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`%s expected to be basetypes.NumberValue, was: %T`, multiplierAttributeName, multiplierAttribute))
	}

	randomizationFactorAttribute, ok := attributes[randomizationFactorAttributeName]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			fmt.Sprintf(`%s is missing from object`, randomizationFactorAttributeName))

		return nil, diags
	}

	randomizationFactorVal, ok := randomizationFactorAttribute.(basetypes.NumberValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`%s expected to be basetypes.NumberValue, was: %T`, randomizationFactorAttributeName, randomizationFactorAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return RetryValue{
		ErrorMessageRegex:   errorMessageRegexVal,
		IntervalSeconds:     intervalSecondsVal,
		MaxIntervalSeconds:  maxIntervalSecondsVal,
		Multiplier:          multiplierVal,
		RandomizationFactor: randomizationFactorVal,
		state:               attr.ValueStateKnown,
	}, diags
}

func NewRetryValueNull() RetryValue {
	return RetryValue{
		state: attr.ValueStateNull,
	}
}

func NewRetryValueUnknown() RetryValue {
	return RetryValue{
		state: attr.ValueStateUnknown,
	}
}

func NewRetryValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (RetryValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing RetryableErrorsValue Attribute Value",
				"While creating a RetryableErrorsValue value, a missing attribute value was detected. "+
					"A RetryableErrorsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("RetryableErrorsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid RetryableErrorsValue Attribute Type",
				"While creating a RetryableErrorsValue value, an invalid attribute value was detected. "+
					"A RetryableErrorsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("RetryableErrorsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("RetryableErrorsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra RetryableErrorsValue Attribute Value",
				"While creating a RetryableErrorsValue value, an extra attribute value was detected. "+
					"A RetryableErrorsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra RetryableErrorsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewRetryValueUnknown(), diags
	}

	errorMessageRegexAttribute, ok := attributes[errorMessageRegexAttributeName]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`error_message_regex is missing from object`)

		return NewRetryValueUnknown(), diags
	}

	errorMessageRegexVal, ok := errorMessageRegexAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`error_message_regex expected to be basetypes.ListValue, was: %T`, errorMessageRegexAttribute))
	}

	intervalSecondsAttribute, ok := attributes[intervalSecondsAttributeName]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`interval_seconds is missing from object`)

		return NewRetryValueUnknown(), diags
	}

	intervalSecondsVal, ok := intervalSecondsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`interval_seconds expected to be basetypes.Int64Value, was: %T`, intervalSecondsAttribute))
	}

	maxIntervalSecondsAttribute, ok := attributes[maxIntervalSecondsAttributeName]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`max_interval_seconds is missing from object`)

		return NewRetryValueUnknown(), diags
	}

	maxIntervalSecondsVal, ok := maxIntervalSecondsAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`max_interval_seconds expected to be basetypes.Int64Value, was: %T`, maxIntervalSecondsAttribute))
	}

	multiplierAttribute, ok := attributes[multiplierAttributeName]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`multiplier is missing from object`)

		return NewRetryValueUnknown(), diags
	}

	multiplierVal, ok := multiplierAttribute.(basetypes.NumberValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`multiplier expected to be basetypes.NumberValue, was: %T`, multiplierAttribute))
	}

	randomizationFactorAttribute, ok := attributes[randomizationFactorAttributeName]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`randomization_factor is missing from object`)

		return NewRetryValueUnknown(), diags
	}

	randomizationFactorVal, ok := randomizationFactorAttribute.(basetypes.NumberValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`randomization_factor expected to be basetypes.NumberValue, was: %T`, randomizationFactorAttribute))
	}

	if diags.HasError() {
		return NewRetryValueUnknown(), diags
	}

	return RetryValue{
		ErrorMessageRegex:   errorMessageRegexVal,
		IntervalSeconds:     intervalSecondsVal,
		MaxIntervalSeconds:  maxIntervalSecondsVal,
		Multiplier:          multiplierVal,
		RandomizationFactor: randomizationFactorVal,
		state:               attr.ValueStateKnown,
	}, diags
}

func NewRetryableErrorsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) RetryValue {
	object, diags := NewRetryValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewRetryableErrorsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t RetryType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewRetryValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewRetryValueUnknown(), nil
	}

	if in.IsNull() {
		return NewRetryValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewRetryableErrorsValueMust(RetryValue{}.AttributeTypes(ctx), attributes), nil
}

type RetryValue struct {
	ErrorMessageRegex   types.List   `tfsdk:"error_message_regex"`
	IntervalSeconds     types.Int64  `tfsdk:"interval_seconds"`
	MaxIntervalSeconds  types.Int64  `tfsdk:"max_interval_seconds"`
	Multiplier          types.Number `tfsdk:"multiplier"`
	RandomizationFactor types.Number `tfsdk:"randomization_factor"`
	state               attr.ValueState
}

func (v RetryValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 5)

	var val tftypes.Value
	var err error

	attrTypes[errorMessageRegexAttributeName] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes[intervalSecondsAttributeName] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes[maxIntervalSecondsAttributeName] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes[multiplierAttributeName] = basetypes.NumberType{}.TerraformType(ctx)
	attrTypes[randomizationFactorAttributeName] = basetypes.NumberType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 5)

		val, err = v.ErrorMessageRegex.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals[errorMessageRegexAttributeName] = val

		val, err = v.IntervalSeconds.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals[intervalSecondsAttributeName] = val

		val, err = v.MaxIntervalSeconds.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals[maxIntervalSecondsAttributeName] = val

		val, err = v.Multiplier.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals[multiplierAttributeName] = val

		val, err = v.RandomizationFactor.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals[randomizationFactorAttributeName] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v RetryValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v RetryValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v RetryValue) String() string {
	return "retry.RetryableErrorsValue"
}

func (v RetryValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	errorMessageRegexVal, d := types.ListValue(types.StringType, v.ErrorMessageRegex.Elements())

	diags.Append(d...)

	if d.HasError() {
		return types.ObjectUnknown(map[string]attr.Type{
			errorMessageRegexAttributeName: basetypes.ListType{
				ElemType: types.StringType,
			},
			intervalSecondsAttributeName:     basetypes.Int64Type{},
			maxIntervalSecondsAttributeName:  basetypes.Int64Type{},
			multiplierAttributeName:          basetypes.NumberType{},
			randomizationFactorAttributeName: basetypes.NumberType{},
		}), diags
	}

	attributeTypes := map[string]attr.Type{
		errorMessageRegexAttributeName: basetypes.ListType{
			ElemType: types.StringType,
		},
		intervalSecondsAttributeName:     basetypes.Int64Type{},
		maxIntervalSecondsAttributeName:  basetypes.Int64Type{},
		multiplierAttributeName:          basetypes.NumberType{},
		randomizationFactorAttributeName: basetypes.NumberType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			errorMessageRegexAttributeName:   errorMessageRegexVal,
			intervalSecondsAttributeName:     v.IntervalSeconds,
			maxIntervalSecondsAttributeName:  v.MaxIntervalSeconds,
			multiplierAttributeName:          v.Multiplier,
			randomizationFactorAttributeName: v.RandomizationFactor,
		})

	return objVal, diags
}

func (v RetryValue) Equal(o attr.Value) bool {
	other, ok := o.(RetryValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.ErrorMessageRegex.Equal(other.ErrorMessageRegex) {
		return false
	}

	if !v.IntervalSeconds.Equal(other.IntervalSeconds) {
		return false
	}

	if !v.MaxIntervalSeconds.Equal(other.MaxIntervalSeconds) {
		return false
	}

	if !v.Multiplier.Equal(other.Multiplier) {
		return false
	}

	if !v.RandomizationFactor.Equal(other.RandomizationFactor) {
		return false
	}

	return true
}

func (v RetryValue) Type(ctx context.Context) attr.Type {
	return basetypes.ObjectType{
		AttrTypes: v.AttributeTypes(ctx),
	}
}

func (v RetryValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		errorMessageRegexAttributeName: basetypes.ListType{
			ElemType: types.StringType,
		},
		intervalSecondsAttributeName:     basetypes.Int64Type{},
		maxIntervalSecondsAttributeName:  basetypes.Int64Type{},
		multiplierAttributeName:          basetypes.NumberType{},
		randomizationFactorAttributeName: basetypes.NumberType{},
	}
}

func (v RetryValue) GetErrorMessages() []string {
	if v.IsNull() {
		return nil
	}
	if v.IsUnknown() {
		return nil
	}
	res := make([]string, len(v.ErrorMessageRegex.Elements()))
	for i, elem := range v.ErrorMessageRegex.Elements() {
		res[i] = elem.(types.String).ValueString()
	}
	return res
}

func (v RetryValue) GetErrorMessagesRegex() []regexp.Regexp {
	msgs := v.GetErrorMessages()
	if msgs == nil {
		return nil
	}
	res := make([]regexp.Regexp, len(msgs))
	for i, msg := range msgs {
		res[i] = *regexp.MustCompile(msg)
	}
	return res
}

func (v RetryValue) GetIntervalSeconds() int {
	return v.getInt64AttrValue(intervalSecondsAttributeName)
}

func (v RetryValue) GetIntervalSecondsAsDuration() time.Duration {
	return time.Duration(v.IntervalSeconds.ValueInt64()) * time.Second
}

func (v RetryValue) GetMaxIntervalSeconds() int {
	return v.getInt64AttrValue(maxIntervalSecondsAttributeName)
}

func (v RetryValue) GetMaxIntervalSecondsAsDuration() time.Duration {
	return time.Duration(v.MaxIntervalSeconds.ValueInt64()) * time.Second
}

func (v RetryValue) GetMultiplier() float64 {
	return v.getNumberAttrValue(multiplierAttributeName)
}

func (v RetryValue) GetRandomizationFactor() float64 {
	return v.getNumberAttrValue(randomizationFactorAttributeName)
}

func (v RetryValue) getNumberAttrValue(name string) float64 {
	switch name {
	case multiplierAttributeName:
		return bigFloat2Float64(v.Multiplier.ValueBigFloat())
	case randomizationFactorAttributeName:
		return bigFloat2Float64(v.RandomizationFactor.ValueBigFloat())
	default:
		return 0
	}
}

func bigFloat2Float64(bf *big.Float) float64 {
	f, _ := bf.Float64()
	return f
}

func (v RetryValue) getInt64AttrValue(name string) int {
	switch name {
	case intervalSecondsAttributeName:
		return int(v.IntervalSeconds.ValueInt64())
	case maxIntervalSecondsAttributeName:
		return int(v.MaxIntervalSeconds.ValueInt64())
	default:
		return 0
	}
}

func (v RetryValue) AddDefaultValuesIfUnknownOrNull() RetryValue {
	if v.IntervalSeconds.IsUnknown() || v.IntervalSeconds.IsNull() {
		v.IntervalSeconds = basetypes.NewInt64Value(defaultIntervalSeconds)
	}
	if v.MaxIntervalSeconds.IsUnknown() || v.MaxIntervalSeconds.IsNull() {
		v.MaxIntervalSeconds = basetypes.NewInt64Value(defaultMaxIntervalSeconds)
	}
	if v.Multiplier.IsUnknown() || v.Multiplier.IsNull() {
		v.Multiplier = basetypes.NewNumberValue(big.NewFloat(defaultMultiplier))
	}
	if v.RandomizationFactor.IsUnknown() || v.RandomizationFactor.IsNull() {
		v.RandomizationFactor = basetypes.NewNumberValue(big.NewFloat(defaultRandomizationFactor))
	}
	return v
}
