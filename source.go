package whub

import "time"

type Source interface {
	Pull(timeout time.Duration) Message
}
