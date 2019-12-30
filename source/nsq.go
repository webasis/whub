package source

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/webasis/whub"
)

type NATSSource struct {
	conn *nats.Conn
	r    <-chan whub.Message
}

func NATS(conn *nats.Conn, subject string) *NATSSource {
	r := make(chan whub.Message, 16)

	src := &NATSSource{
		conn: conn,
		r:    r,
	}
	src.conn.Subscribe(subject, func(msg *nats.Msg) {
		m := whub.M()
		err := json.Unmarshal(msg.Data, m)
		if err != nil {
			//TODO log
			return
		}

		r <- m
	})

	return src
}

func (src *NATSSource) Pull(timeout time.Duration) whub.Message {
	if timeout == 0 {
		return <-src.r
	}

	select {
	case m, ok := <-src.r:
		if ok {
			return m
		}
	case <-time.NewTimer(timeout).C:
	}
	return nil
}
