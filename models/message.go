package models
import (
	"log"
	"encoding/json"
)
// Message ...
type Message struct {
	Type		string `json:"type,omitempty"`	// event || subscription || unsubscription || broadcast
	Event		string `json:"event"`
	Data		string `json:"data"`
}
func (msg *Message) Encode() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Println("Message encode:", err)
	}
	return b
}
