package zencoding

import (
	"encoding/base64"
	"fmt"
)

func getString(m map[string]interface{}, key string) (string, error) {
	raw, ok := m[key]
	if !ok {
		return "", fmt.Errorf("Key '%s' not found", key)
	}
	s, isString := raw.(string)
	if !isString {
		return "", fmt.Errorf("Key '%s' was not a string", key)
	}
	return s, nil
}

func getStringPointer(m map[string]interface{}, key string) *string {
	if m[key] == nil {
		return nil
	}
	str := m[key].(string)
	return &str
}

func getStringArray(m map[string]interface{}, key string) []string {
	if m[key] == nil {
		return nil
	}
	raw := m[key].([]interface{})
	arr := make([]string, len(raw))
	for idx, val := range raw {
		arr[idx] = val.(string)
	}
	return arr
}

func getBytes(m map[string]interface{}, key string) []byte {
	if m[key] == nil {
		return nil
	}
	s := m[key].(string)
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}
