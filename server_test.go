package whub

import (
	"net/http"
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	if os.Getenv("test") != "server" {
		return
	}

	s := NewServer()
	go http.ListenAndServe(":8998", s)

	for msg := range s.R {
		if msg.Meta().V("to") == "#all" {
			s.LiteC <- func() {
				for id, _ := range s.Agents {
					msg.Meta().Put("to", id)
					s.Send(msg.Clone())
				}
			}
		}
	}
}
