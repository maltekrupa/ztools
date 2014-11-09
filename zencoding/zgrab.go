package zencoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

type EventType interface {
	TypeName() string
	GetEmptyInstance() EventData
	MarshalJSON() ([]byte, error)
}

var typeNameToTypeMap map[string]EventType

func init() {
	typeNameToTypeMap = make(map[string]EventType)
}

func RegisterEventType(t EventType) {
	name := t.TypeName()
	if _, exists := typeNameToTypeMap[name]; exists {
		panic("Duplicate type name " + name)
	}
	typeNameToTypeMap[name] = t
}

func EventTypeFromName(name string) (EventType, error) {
	t, ok := typeNameToTypeMap[name]
	if !ok {
		return t, fmt.Errorf("Unknown type name %s", name)
	}
	return t, nil
}

type encodedGrab struct {
	Host   string            `json:"host"`
	Domain *string           `json:"domain"`
	Time   string            `json:"time"`
	Log    []ConnectionEvent `json:"log"`
}

type encodedConnectionEvent struct {
	Type  EventType `json:"type"`
	Data  EventData `json:"data"`
	Error *string   `json:"error"`
}

type partialConnectionEvent struct {
	Data  EventData `json:"data"`
	Error *string   `json:"error"`
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
	ece := new(partialConnectionEvent)
	tn := struct {
		TypeName string `json:"type"`
	}{}
	if err := json.Unmarshal(b, &tn); err != nil {
		return err
	}
	t, typeErr := EventTypeFromName(tn.TypeName)
	if typeErr != nil {
		return typeErr
	}
	ece.Data = t.GetEmptyInstance()
	if err := json.Unmarshal(b, &ece); err != nil {
		log.Print(err)
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
