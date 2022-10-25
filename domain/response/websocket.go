package response

import "github.com/goccy/go-json"

type WSResponse struct {
	RequestID uint64                 `json:"request_id,omitempty"`
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// FromBytes converts []byte to Input
func (w *WSResponse) FromBytes(buffer []byte) error {
	return json.Unmarshal(buffer, w)
}

// ToBytes converts Input to []byte
func (w *WSResponse) ToBytes() []byte {
	buffer, _ := json.Marshal(w)
	return buffer
}
