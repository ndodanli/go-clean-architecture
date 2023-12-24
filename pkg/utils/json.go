package utils

import "encoding/json"

func ToJson(i interface{}) string {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
