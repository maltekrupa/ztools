package zencoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

type EventData interface {
	GetType() EventType
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type ConnectionEvent struct {
	Data  EventData
	Error error
}

type Grab struct {
	Host   net.IP            `json:"host"`
	Domain string            `json:"domain"`
	Time   time.Time         `json:"timestamp"`
	Log    []ConnectionEvent `json:"log"`
}

type EventType uint8

var typeNameToTypeMap map[string]EventType
var typeToTypeNameMap map[EventType]string

func init() {
	typeNameToTypeMap = make(map[string]EventType)
	typeToTypeNameMap = make(map[EventType]string)

	typeNameToTypeMap[CONNECTION_EVENT_CONNECT_NAME] = CONNECTION_EVENT_CONNECT
	typeNameToTypeMap[CONNECTION_EVENT_TLS_NAME] = CONNECTION_EVENT_TLS

	typeToTypeNameMap[CONNECTION_EVENT_CONNECT] = CONNECTION_EVENT_CONNECT_NAME
	typeToTypeNameMap[CONNECTION_EVENT_TLS] = CONNECTION_EVENT_TLS_NAME
}

func EventTypeFromName(name string) (EventType, error) {
	t, ok := typeNameToTypeMap[name]
	if !ok {
		return t, fmt.Errorf("Unknown type name %s", name)
	}
	return t, nil
}

func (t EventType) TypeName() string {
	return typeToTypeNameMap[t]
}

func (t EventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.TypeName())
}

func (t *EventType) UnmarshalJSON(b []byte) error {
	var typeName string
	if unmarshalErr := json.Unmarshal(b, &typeName); unmarshalErr != nil {
		return unmarshalErr
	}
	eventType, eventTypeErr := EventTypeFromName(typeName)
	if eventTypeErr != nil {
		return eventTypeErr
	}
	*t = eventType
	return nil
}

func (ce *ConnectionEvent) MarshalJSON() ([]byte, error) {
	t := ce.Data.GetType()
	var esp *string
	if ce.Error != nil {
		es := ce.Error.Error()
		esp = &es
	}
	obj := encodedConnectionEvent{
		Type:  t,
		Data:  ce.Data,
		Error: esp,
	}
	return json.Marshal(obj)
}

func (ce *ConnectionEvent) UnmarshalJSON(b []byte) error {
	ece := new(encodedConnectionEvent)
	t := struct {
		Type EventType `json:"type"`
	}{}
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	switch t.Type {
	case CONNECTION_EVENT_TLS:
		ece.Data = new(ServerHandshake)
	default:
		return fmt.Errorf("Unknown event type: %s", t.Type.TypeName())
	}
	if err := json.Unmarshal(b, &ece); err != nil {
		return err
	}
	ce.Data = ece.Data
	if ece.Error != nil {
		ce.Error = errors.New(*ece.Error)
	}
	return nil
}

func (g *Grab) MarshalJSON() ([]byte, error) {
	var domainPtr *string
	if g.Domain != "" {
		domainPtr = &g.Domain
	}
	time := g.Time.Format(time.RFC3339)
	obj := encodedGrab{
		Host:   g.Host.String(),
		Domain: domainPtr,
		Time:   time,
		Log:    g.Log,
	}
	return json.Marshal(obj)
}

func (g *Grab) UnmarshalJSON(b []byte) error {
	eg := new(encodedGrab)
	err := json.Unmarshal(b, eg)
	if err != nil {
		return err
	}
	g.Host = net.ParseIP(eg.Host)
	if eg.Domain != nil {
		g.Domain = *eg.Domain
	}
	if g.Time, err = time.Parse(time.RFC3339, eg.Time); err != nil {
		return err
	}
	g.Log = eg.Log
	return nil
}

const (
	CONNECTION_EVENT_CONNECT      EventType = 0
	CONNECTION_EVENT_CONNECT_NAME           = "connect"
	CONNECTION_EVENT_TLS          EventType = 1
	CONNECTION_EVENT_TLS_NAME               = "tls_handshake"
)
