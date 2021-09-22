package utils

import (
	"encoding/json"
	"fmt"
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

func GetMergedJson(old interface{}, new interface{}) interface{} {
	switch old.(type) {
	case map[string]interface{}:
		switch new.(type) {
		case map[string]interface{}:
			oldMap := old.(map[string]interface{})
			newMap := new.(map[string]interface{})
			res := make(map[string]interface{})
			for key, oldValue := range oldMap {
				if newMap[key] != nil {
					res[key] = GetMergedJson(oldValue, newMap[key])
				} else {
					res[key] = oldValue
				}
			}
			for key, newValue := range newMap {
				if res[key] == nil {
					res[key] = newValue
				}
			}
			return res
		default:
			return new
		}
	default:
		return new
	}
}

func GetUpdatedJson(old interface{}, new interface{}) interface{} {
	switch old.(type) {
	case map[string]interface{}:
		switch new.(type) {
		case map[string]interface{}:
			oldMap := old.(map[string]interface{})
			newMap := new.(map[string]interface{})
			res := make(map[string]interface{})
			for key, oldValue := range oldMap {
				if newMap[key] != nil {
					res[key] = GetUpdatedJson(oldValue, newMap[key])
				}
			}
			return res
		default:
			return new
		}
	case []interface{}:
		switch new.(type) {
		case []interface{}:
			oldArr := old.([]interface{})
			newArr := new.([]interface{})
			if len(oldArr) != len(newArr) {
				return newArr
			}
			res := make([]interface{}, 0)
			for index := range oldArr {
				res = append(res, GetUpdatedJson(oldArr[index], newArr[index]))
			}
			return res
		default:
			return new
		}
	default:
		return new
	}
}

func GetRemovedJson(old interface{}, new interface{}) interface{} {
	switch old.(type) {
	case map[string]interface{}:
		switch new.(type) {
		case map[string]interface{}:
			oldMap := old.(map[string]interface{})
			newMap := new.(map[string]interface{})
			res := make(map[string]interface{})
			for key, oldValue := range oldMap {
				if newMap[key] != nil {
					res[key] = GetRemovedJson(oldValue, newMap[key])
				} else {
					res[key] = oldValue
				}
			}
			return res
		default:
			return nil
		}
	default:
		return nil
	}
}

func GetIgnoredJson(old interface{}, ignoredProperties []string) interface{} {
	switch old.(type) {
	case map[string]interface{}:
		oldMap := old.(map[string]interface{})
		res := make(map[string]interface{})
		for key, value := range oldMap {
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
		oldArr := old.([]interface{})
		res := make([]interface{}, 0)
		for index := range oldArr {
			res = append(res, GetIgnoredJson(oldArr[index], ignoredProperties))
		}
		return res
	default:
		return old
	}
}
