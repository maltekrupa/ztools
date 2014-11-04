package zencoding

import (
	"encoding/json"
	"time"
)

type EventData interface {
	GetType() EventType
	UnpackMap(map[string]interface{}) error
	MarshalJSON() ([]byte, error)
}

type ConnectionEvent struct {
	Data  EventData
	Error error
}

type Grab struct {
	Host   string            `json:"host"`
	Domain string            `json:"domain"`
	Time   time.Time         `json:"timestamp"`
	Log    []ConnectionEvent `json:"log"`
}

type EventType interface {
	String() string
}

func (ce *ConnectionEvent) MarshalJSON() ([]byte, error) {
	t := ce.Data.GetType().String()
	var esp *string
	if ce.Error != nil {
		es := ce.Error.Error()
		esp = &es
	}
	obj := struct {
		Type  string    `json:"type"`
		Data  EventData `json:"data"`
		Error *string   `json:"error"`
	}{
		Type:  t,
		Data:  ce.Data,
		Error: esp,
	}
	return json.Marshal(obj)
}

func (g *Grab) MarshalJSON() ([]byte, error) {
	var domainPtr *string
	if g.Domain != "" {
		domainPtr = &g.Domain
	}
	time := g.Time.UTC().Format(time.RFC3339)
	obj := struct {
		Host   string            `json:"host"`
		Domain *string           `json:"domain"`
		Time   string            `json:"time"`
		Log    []ConnectionEvent `json:"log"`
	}{
		Host:   g.Host,
		Domain: domainPtr,
		Time:   time,
		Log:    g.Log,
	}
	return json.Marshal(obj)
}
