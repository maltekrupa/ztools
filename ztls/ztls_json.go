package ztls

import (
	"encoding/json"
	"ztools/zencoding"
)

type encodedHandshake struct {
	Hello        *ServerHello        `json:"server_hello"`
	Certificates *ServerCertificates `json:"server_certificates"`
	KeyExchange  *ServerKeyExchange  `json:"server_key_exchange"`
	Finished     *ServerFinished     `json:"server_finished"`
}

// Implements zencoding.EventType
type EventTypeTLS uint8

func (t EventTypeTLS) TypeName() string {
	return CONNECTION_EVENT_TLS_NAME
}

func (t EventTypeTLS) GetEmptyInstance() zencoding.EventData {
	return new(ServerHandshake)
}

func (t EventTypeTLS) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.TypeName())
}

// GetType always returns the TLS Handshake type
func (hs *ServerHandshake) GetType() zencoding.EventType {
	var t EventTypeTLS
	return t
}

// MarshalJSON implements the json.Marshaler interface
func (hs *ServerHandshake) MarshalJSON() ([]byte, error) {
	// Prevent infinite recursion
	obj := encodedHandshake{
		Hello:        hs.ServerHello,
		Certificates: hs.ServerCertificates,
		KeyExchange:  hs.ServerKeyExchange,
		Finished:     hs.ServerFinished,
	}
	return json.Marshal(obj)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (hs *ServerHandshake) UnmarshalJSON(b []byte) error {
	obj := encodedHandshake{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	hs.ServerHello = obj.Hello
	hs.ServerCertificates = obj.Certificates
	hs.ServerKeyExchange = obj.KeyExchange
	hs.ServerFinished = obj.Finished
	return nil
}

func init() {
	var t EventTypeTLS
	zencoding.RegisterEventType(t)
}

const (
	CONNECTION_EVENT_TLS_NAME = "tls_handshake"
)
