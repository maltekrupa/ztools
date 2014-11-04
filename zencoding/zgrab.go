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
	UnpackMap(map[string]interface{}) error
	MarshalJSON() ([]byte, error)
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

func (ce *ConnectionEvent) UnpackMap(raw map[string]interface{}) error {
	eventType, err := getString(raw, "type")
	if err != nil {
		return err
	}
	eventError := getStringPointer(raw, "error")
	if eventError == nil {
		ce.Error = errors.New(*eventError)
	} else {
		ce.Data = eventDataFromTypeName(eventType)
		if ce.Data == nil {
			return fmt.Errorf("Unknown event type %s", eventType)
		}
		rawData, ok := raw["data"].(map[string]interface{})
		if !ok {
			return errors.New("Invalid \"data\" field")
		}
		ce.Data.UnpackMap(rawData)
	}
	return nil
}

func (g *Grab) MarshalJSON() ([]byte, error) {
	var domainPtr *string
	if g.Domain != "" {
		domainPtr = &g.Domain
	}
	time := g.Time.Format(time.RFC3339)
	obj := struct {
		Host   string            `json:"host"`
		Domain *string           `json:"domain"`
		Time   string            `json:"time"`
		Log    []ConnectionEvent `json:"log"`
	}{
		Host:   g.Host.String(),
		Domain: domainPtr,
		Time:   time,
		Log:    g.Log,
	}
	return json.Marshal(obj)
}

func (g *Grab) UnpackMap(raw map[string]interface{}) error {
	var s string
	var sp *string
	var err error
	// Read the IP as the "host" field
	if s, err = getString(raw, "host"); err != nil {
		return err
	}
	if g.Host = net.ParseIP(s); g.Host == nil {
		return fmt.Errorf("Invalid host: %s is not an IP", s)
	}
	// Read the domain
	if sp = getStringPointer(raw, "domain"); sp != nil {
		g.Domain = *sp
	}
	// Read the timestamp
	if s, err = getString(raw, "time"); err != nil {
		return err
	}
	if g.Time, err = time.Parse(time.RFC3339, s); err != nil {
		return err
	}
	// Read each element of the log
	log := raw["log"].([]interface{})
	g.Log = make([]ConnectionEvent, len(log))
	for idx, val := range log {
		rawEvent, ok := val.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Invalid log entry at index %d", idx)
		}
		if err = g.Log[idx].UnpackMap(rawEvent); err != nil {
			return fmt.Errorf("Invalid log entry at index %d: %s", idx, err.Error())
		}
	}
	return nil
}

func eventDataFromTypeName(name string) EventData {
	switch name {
	case CONNECTION_EVENT_TLS_NAME:
		return new(ServerHandshake)
	default:
		return nil
	}
}

const (
	CONNECTION_EVENT_TLS_NAME = "tls_handshake"
)
