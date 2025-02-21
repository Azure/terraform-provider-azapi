package skip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type SkipTestType1 struct {
	SkipOnUpdate     types.String `skip_on:"update"`
	NoSkip           types.String // no skip_on tag
	SkipOnCreateRead types.String `tfsdk:"field_3" skip_on:"create,read"`
}

func TestCanSkipExternalRequest(t *testing.T) {
	testCases := []struct {
		desc      string
		a         SkipTestType1
		b         SkipTestType1
		operation string
		result    bool
	}{
		{
			desc: "skip on update",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1_updated"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			operation: "update",
			result:    true,
		},
		{
			desc: "do not skip on update, changes to untagged fields",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2_updated"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			operation: "update",
			result:    false,
		},
		{
			desc: "skip on create",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3_updated"),
			},
			operation: "create",
			result:    true,
		},
		{
			desc: "do not skip on create, changes to non-skippable fields",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2_updated"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			operation: "create",
			result:    false,
		},
		{
			desc: "skip on read",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3_updated"),
			},
			operation: "read",
			result:    true,
		},
		{
			desc: "do not skip on read, changes to non-skippable fields",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2_updated"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			operation: "read",
			result:    false,
		},
		{
			desc: "skip on update but no changes",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			operation: "update",
			result:    true,
		},
		{
			desc: "skip on update, changes to read skippable field",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3_updated"),
			},
			operation: "update",
			result:    false,
		},
		{
			desc: "skip on update, unknown value",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringUnknown(),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringUnknown(),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			operation: "update",
			result:    true,
		},
		{
			desc: "skip on update, null value",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringNull(),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringNull(),
				NoSkip:           basetypes.NewStringValue("value2"),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			operation: "update",
			result:    true,
		},
		{
			desc: "no skip, null and unknown value differ",
			a: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringUnknown(),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			b: SkipTestType1{
				SkipOnUpdate:     basetypes.NewStringValue("value1"),
				NoSkip:           basetypes.NewStringNull(),
				SkipOnCreateRead: basetypes.NewStringValue("value3"),
			},
			operation: "update",
			result:    false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if CanSkipExternalRequest(tC.a, tC.b, tC.operation) != tC.result {
				t.Errorf("Expected %v, got %v", tC.result, !tC.result)
			}
		})
	}
}

func TestCanSkipExternalRequestInvalidInputs(t *testing.T) {
	// Test with non-struct values
	if CanSkipExternalRequest(1, 2, "update") {
		t.Errorf("Expected false, got true")
	}

}
