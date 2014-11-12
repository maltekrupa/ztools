package zencoding

import (
	"encoding/json"
	"fmt"
)

type EventData interface {
	GetType() EventType
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type EventType struct {
	TypeName         string
	GetEmptyInstance func() EventData
}

// MarshalJSON implements the json.Marshaler interface
func (e EventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.TypeName)
}

var typeNameToTypeMap map[string]EventType

func init() {
	typeNameToTypeMap = make(map[string]EventType)
	RegisterEventType(TLSHandshakeEventType)
}

func RegisterEventType(t EventType) {
	name := t.TypeName
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
