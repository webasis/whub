package whub

type Handler interface {
	On(m Message)
}

type HandleFunc func(m Message)

func (fn HandleFunc) On(m Message) { fn(m) }

type Router interface {
	Route(m Message) Handler
}

type RouteFunc func(m Message) Handler

func (fn RouteFunc) Route(m Message) Handler { return fn(m) }
