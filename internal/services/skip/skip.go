package skip

import (
	"reflect"
	"slices"
	"strings"
)

// CanSkipExternalRequest checks if the external request can be skipped based on the plan and state.
// Two of the same objects are supplied as parameters, together with the operation that is being performed.
// The function uses the `skip_on` struct tag to determine if the field should be skipped.
// The value of the `skip_on` tag is a comma-separated list of operations that mean that changes to this field value do not require an external request and are in state only.
// The function will return true if the external request can be skipped, false otherwise.
func CanSkipExternalRequest[T any](a, b T, operation string) bool {
	valA := reflect.ValueOf(a)
	valB := reflect.ValueOf(b)

	// Since we are using generics, we know that the types of a and b are the same.
	// Therefore we can check the type of a to determine if it is a struct.
	if valA.Kind() != reflect.Struct {
		return false
	}

	typeOfA := valA.Type()
	// iterate over all fields of the struct
	for i := 0; i < typeOfA.NumField(); i++ {
		field := typeOfA.Field(i)
		// Check if the field has the skip_on tag
		// If it doesn't we need to compare the valued as we cannot determine if the field should be skipped.
		// If the field has the skip_on tag, we can check if the operation is in the list of operations that should be skipped.
		tag := field.Tag.Get("skip_on")
		if tag != "" {
			// Split the tag values by comma and check if the operation is in the list.
			// If the operation is in the list, then this field represents a change in state only
			// and does not require an external request to be made.
			// Therefore we can skip tp the next field.
			tagValues := strings.Split(tag, ",")
			if slices.Contains(tagValues, operation) {
				continue
			}
		}

		// If we get here then we need to compare the field values.
		// By now we have determined that the struct fields do not have a valid skip value for this operation.
		// Therefore if the field values are not equal, then the external request cannot be skipped.
		fieldA := valA.Field(i)
		fieldB := valB.Field(i)

		if !fieldA.IsValid() || !fieldB.IsValid() {
			return false
		}
		if !reflect.DeepEqual(fieldA.Interface(), fieldB.Interface()) {
			return false
		}
	}
	return true
}
