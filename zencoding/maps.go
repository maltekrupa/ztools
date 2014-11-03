package zencoding

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
