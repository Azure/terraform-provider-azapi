package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func NormalizeJson(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}
	var j interface{}

	if err := json.Unmarshal([]byte(jsonString.(string)), &j); err != nil {
		return fmt.Sprintf("Error parsing JSON: %+v", err)
	}
	b, _ := json.Marshal(j)
	return string(b)
}

// GetMergedJson is used to merge object old and new, if overlaps, use new value
func GetMergedJson(old interface{}, new interface{}) interface{} {
	if new == nil {
		return new
	}
	switch oldValue := old.(type) {
	case map[string]interface{}:
		if newMap, ok := new.(map[string]interface{}); ok {
			res := make(map[string]interface{})
			for key, value := range oldValue {
				if _, ok := newMap[key]; ok {
					res[key] = GetMergedJson(value, newMap[key])
				} else {
					res[key] = value
				}
			}
			for key, newValue := range newMap {
				if res[key] == nil {
					res[key] = newValue
				}
			}
			return res
		}
	case []interface{}:
		if newArr, ok := new.([]interface{}); ok {
			if len(oldValue) != len(newArr) {
				return newArr
			}
			res := make([]interface{}, 0)
			for index := range oldValue {
				res = append(res, GetMergedJson(oldValue[index], newArr[index]))
			}
			return res
		}
	}
	return new
}

type UpdateJsonOption struct {
	IgnoreCasing          bool
	IgnoreMissingProperty bool
}

// GetUpdatedJson is used to get an updated object which has same schema as old, but with new value
func GetUpdatedJson(old interface{}, new interface{}, option UpdateJsonOption) interface{} {
	switch oldValue := old.(type) {
	case map[string]interface{}:
		if newMap, ok := new.(map[string]interface{}); ok {
			res := make(map[string]interface{})
			for key, value := range oldValue {
				switch {
				case newMap[key] != nil:
					res[key] = GetUpdatedJson(value, newMap[key], option)
				case option.IgnoreMissingProperty || isZeroValue(value):
					res[key] = value
				}
			}
			return res
		}
	case []interface{}:
		if newArr, ok := new.([]interface{}); ok {
			if len(oldValue) != len(newArr) {
				return newArr
			}
			res := make([]interface{}, 0)
			for index := range oldValue {
				res = append(res, GetUpdatedJson(oldValue[index], newArr[index], option))
			}
			return res
		}
	case string:
		if newStr, ok := new.(string); ok {
			if option.IgnoreCasing && strings.EqualFold(oldValue, newStr) {
				return oldValue
			}
			if option.IgnoreMissingProperty && regexp.MustCompile(`^\**$`).MatchString(newStr) {
				return oldValue
			}
		}
	}
	return new
}

// GetRemovedJson is used to get an object which is remove properties defined in new from old
func GetRemovedJson(old interface{}, new interface{}) interface{} {
	switch oldValue := old.(type) {
	case map[string]interface{}:
		if newMap, ok := new.(map[string]interface{}); ok {
			res := make(map[string]interface{})
			for key, value := range oldValue {
				if newMap[key] != nil {
					value := GetRemovedJson(value, newMap[key])
					if value != nil {
						res[key] = value
					}
				} else {
					res[key] = value
				}
			}
			if len(res) == 0 {
				return nil
			}
			return res
		}
	case []interface{}:
		if newArr, ok := new.([]interface{}); ok {
			if len(oldValue) != len(newArr) {
				return nil
			}
			res := make([]interface{}, 0)
			for index := range oldValue {
				res = append(res, GetRemovedJson(oldValue[index], newArr[index]))
			}
			if len(res) == 0 {
				return nil
			}
			return res
		}
	}
	return nil
}

// GetIgnoredJson is used to remove properties which is in the list called ignoredProperties
func GetIgnoredJson(old interface{}, ignoredProperties []string) interface{} {
	switch oldValue := old.(type) {
	case map[string]interface{}:
		res := make(map[string]interface{})
		for key, value := range oldValue {
			found := false
			for _, prop := range ignoredProperties {
				if prop == key {
					found = true
					break
				}
			}
			if !found {
				res[key] = GetIgnoredJson(value, ignoredProperties)
			}
		}
		return res
	case []interface{}:
		res := make([]interface{}, 0)
		for index := range oldValue {
			res = append(res, GetIgnoredJson(oldValue[index], ignoredProperties))
		}
		return res
	default:
		return old
	}
}

// ExtractObject is used to extract object from old for a json path
func ExtractObject(old interface{}, path string) interface{} {
	if len(path) == 0 {
		return old
	}
	if oldMap, ok := old.(map[string]interface{}); ok {
		index := strings.Index(path, ".")
		if index != -1 {
			key := path[0:index]
			result := make(map[string]interface{}, 1)
			value := ExtractObject(oldMap[key], path[index+1:])
			if value == nil {
				return nil
			} else {
				result[key] = value
			}
			return result
		} else {
			if oldMap[path] != nil {
				result := make(map[string]interface{}, 1)
				result[path] = oldMap[path]
				return result
			} else {
				return nil
			}
		}
	}
	return nil
}

// NormalizeObject is used to remove customized type and replaced with builtin type
func NormalizeObject(input interface{}) interface{} {
	jsonString, _ := json.Marshal(input)
	var output interface{}
	_ = json.Unmarshal(jsonString, &output)
	return output
}

func isZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case map[string]interface{}:
		return len(v) == 0
	case []interface{}:
		return len(v) == 0
	case string:
		return len(v) == 0
	case int, int32, int64, float32, float64:
		return v == 0
	case bool:
		return !v
	}
	return false
}
