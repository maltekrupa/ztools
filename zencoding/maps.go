package zencoding

import "encoding/base64"

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
