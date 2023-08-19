package helper

import (
	"encoding/json"
	"io"
)

// helper to first read all from resp and then unmarshal (better for invalid json)
func UnmarshalResp(resp io.Reader, v interface{}) error {
	body, err := io.ReadAll(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
