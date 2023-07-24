package utils

import (
	"encoding/json"
)

func StructToMap(s any) map[string]any {
	var m map[string]any
	inrec, _ := json.Marshal(s)
	json.Unmarshal(inrec, &m)

	return m
}
