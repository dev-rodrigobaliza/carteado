package request

import "github.com/goccy/go-json"

type WSRequest struct {
	RequestID uint64                 `json:"request_id"`
	Service   string                 `json:"service"`
	Resource  string                 `json:"resource"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// FromBytes converts []byte to Input
func (w *WSRequest) FromBytes(buffer []byte) error {
	return json.Unmarshal(buffer, w)
}

// ToBytes converts Input to []byte
func (w *WSRequest) ToBytes() []byte {
	buffer, _ := json.Marshal(w)
	return buffer
}
