package whub

import "time"

const (
	// Time allowed to write the file to the client.
	WriteWait = 10 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	PingPeriod = 60 * time.Second
)
