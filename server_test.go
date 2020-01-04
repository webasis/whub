package whub

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	if os.Getenv("test") != "server" {
		return
	}

	s := NewServer()
	http.Handle("/ws", s)
	http.HandleFunc("/note", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		m := M()
		m.R("args").
			Put("key", "note").
			Put("value", string(data))
		m.Meta().Put("to", "@set")
		s.R <- m
		w.WriteHeader(201)
	})

	go http.ListenAndServe(":8998", nil)

	kv := make(map[string]string)
	for msg := range s.R {
		fmt.Println("on:", msg)
		switch msg.Meta().V("to") {
		case "@get":
			key := msg.R("args").V("key")
			resp := M()
			resp.R("body").
				Put("key", key).
				Put("value", kv[key])
			resp.Meta().Put("to", msg.Meta().V("from")).
				Put("type", "set")
			s.W <- resp
		case "@set":
			key := msg.R("args").V("key")
			value := msg.R("args").V("value")
			kv[key] = value

			// broadcast
			s.LiteC <- func() {
				m := M()
				m.R("body").
					Put("key", key).
					Put("value", value)
				m.Meta().Put("type", "set")
				for id, _ := range s.Agents {
					m.Meta().Put("to", id)
					s.Send(m.Clone())
				}
			}
		}
	}
}
