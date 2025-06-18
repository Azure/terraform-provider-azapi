package fix

// GetWriteOnlyFix is a function that contains hardcoded logic to fix the write-only properties in the resource body.
func GetWriteOnlyFix(input interface{}) interface{} {
	if input == nil {
		return nil
	}
	input = fixUserAssignedIdentities(input)
	return input
}

// fixUserAssignedIdentities is a helper function to ensure that userAssignedIdentities in the identity object are set to an empty map
func fixUserAssignedIdentities(input interface{}) interface{} {
	if input == nil {
		return nil
	}
	if v, ok := input.(map[string]interface{}); ok && v["identity"] != nil {
		if identity, ok := v["identity"].(map[string]interface{}); ok && identity["userAssignedIdentities"] != nil {
			if userAssignedIdentities, ok := identity["userAssignedIdentities"].(map[string]interface{}); ok {
				for key := range userAssignedIdentities {
					// Ensure that each userAssignedIdentity is an empty map
					userAssignedIdentities[key] = map[string]interface{}{}
				}
			}
			return v
		}
	}
	return input
}
