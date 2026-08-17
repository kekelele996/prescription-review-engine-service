package database

import "encoding/json"

// jsonMarshal JSON 序列化辅助。
func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
