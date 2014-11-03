package zencoding

import "time"

type EventData interface {
}

type ConnectionEvent struct {
	Type  string    `json:"type"`
	Data  EventData `json:"data"`
	Error *string   `json:"error"`
}

type Grab struct {
	Host     string            `json:"host"`
	Domain   *string           `json:"domain"`
	Port     uint16            `json:"port"`
	Protocol string            `json:"protocol"`
	Time     time.Time         `json:"timestamp"`
	Log      []ConnectionEvent `json:"log"`
}
