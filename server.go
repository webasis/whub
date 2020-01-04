package whub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var DefaultUpgrader = websocket.Upgrader{
	ReadBufferSize:  128,
	WriteBufferSize: 128,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Agent struct {
	Conn *websocket.Conn
	Id   string
	C    chan Message // For writing to client.
}

type MiddlewareFunc func(Message) error
type ServerFunc func(*Server)

type Server struct {
	Agents map[string]Agent // map[Id]Agent

	// Important: Both of these should be never closed!!!
	C               chan ServerFunc
	LiteC           chan func()
	R               chan Message // You Should Handle message By your Self
	W               chan Message
	ReadMiddleware  []MiddlewareFunc // MUST NOT USE Server inside middleware
	WriteMiddleware []MiddlewareFunc // MUST NOT USE Server inside middleware

	// READ-ONLY
	Upgrader websocket.Upgrader
	Id       <-chan string
}

func NewServer() *Server {
	Id := make(chan string, 256)
	go func() {
		for i := 1000; ; i++ {
			Id <- fmt.Sprintf("_%d", i)
		}
	}()

	s := &Server{
		Agents:   make(map[string]Agent),
		C:        make(chan ServerFunc, 1024),
		LiteC:    make(chan func(), 1024),
		R:        make(chan Message, 1024*16),
		W:        make(chan Message, 1024*16),
		Upgrader: DefaultUpgrader,
		Id:       Id,
	}
	go s.loop()

	return s
}

func (s *Server) loop() {
	for {
		select {
		case fn, ok := <-s.C:
			if !ok {
				return
			}
			fn(s)
		case fn, ok := <-s.LiteC:
			if !ok {
				return
			}
			fn()
		case msg, ok := <-s.W:
			if !ok {
				return
			}

			s.Send(msg)

		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	agent := s.join(conn)
	defer s.leave(agent.Id)

	go func() {
		pingTicker := time.NewTicker(PingPeriod)
		defer pingTicker.Stop()

		for {
			select {
			case msg, ok := <-agent.C:
				if !ok {
					return
				}

				data, err := json.Marshal(msg)
				if err != nil {
					return
				}
				conn.SetWriteDeadline(time.Now().Add(WriteWait))
				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					conn.Close()
				}
			case <-pingTicker.C:
				err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(WriteWait))
				if err != nil {
					conn.Close()
				}
			}
		}
	}()

	for {
		conn.SetReadDeadline(time.Now().Add(PongWait))
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var msg Message
		if err := json.Unmarshal(p, &msg); err != nil {
			return
		}
		msg.R("meta").Put("from", agent.Id)
		s.R <- msg
	}
}

func (s *Server) Send(msg Message) error {
	agent, ok := s.Agents[msg.R("meta").V("to")]
	if ok {
		for _, middleware := range s.WriteMiddleware {
			err := middleware(msg)
			if err != nil {
				return err
			}
		}

		if msg.R("meta").V("close") != "" {
			agent.Conn.Close()
			return errors.New("meta.close was set")
		}

		select {
		case agent.C <- msg:
			return nil
		default:
			agent.Conn.Close()
			return errors.New("too many message wait to send, so closed the connection")
		}
	}
	return errors.New("not found target")
}

func (s *Server) join(conn *websocket.Conn) Agent {
	c := make(chan Agent, 1)
	s.LiteC <- func() {
		agent := Agent{
			Conn: conn,
			Id:   <-s.Id,
			C:    make(chan Message, 256),
		}

		s.Agents[agent.Id] = agent
		c <- agent
	}
	return <-c
}

func (s *Server) leave(id string) {
	s.LiteC <- func() {
		agent, ok := s.Agents[id]
		if !ok {
			return
		}

		close(agent.C)
		delete(s.Agents, id)
	}
}
