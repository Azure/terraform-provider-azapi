package utils

func GetId(resource interface{}) string {
	if resource == nil {
		return ""
	}
	if resourceMap, ok := resource.(map[string]interface{}); ok {
		if id, ok := resourceMap["id"]; ok {
			return id.(string)
		}
	}
	return ""
}
