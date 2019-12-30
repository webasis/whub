package source

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/webasis/whub"
)

func NATS(conn *nats.Conn, subject string, c chan<- whub.Message) {
	conn.Subscribe(subject, func(msg *nats.Msg) {
		m := whub.M()
		err := json.Unmarshal(msg.Data, &m)
		if err != nil {
			//TODO log
			return
		}

		c <- m
	})
}
