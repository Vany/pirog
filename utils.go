package pirog

import (
	"bytes"
	"encoding/json"
)

func ToJson(in any) string {
	buff := bytes.Buffer{}
	MUST(json.NewEncoder(&buff).Encode(in))
	return buff.String()
}
