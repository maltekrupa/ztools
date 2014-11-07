package zencoding

import "encoding/json"

// ServerHello represents a TLS ServerHello message in a format friendly to the golang JSON library
type ServerHello struct {
	Version            uint16 `json:"version"`
	Random             []byte `json:"random"`
	SessionID          []byte `json:"session_id"`
	CipherSuite        uint16 `json:"cipher_suite"`
	CompressionMethod  uint8  `json:"compression_method"`
	OcspStapling       bool   `json:"ocsp_stapling"`
	TicketSupported    bool   `json:"ticket_supported"`
	HeartbeatSupported bool   `json:"heartbeat_supported"`
}

// ServerCertificates represents a TLS ServerCertificates message in a format friendly to the golang JSON library.
// ValidationError should be non-nil whenever Valid is false.
type ServerCertificates struct {
	Certificates    [][]byte `json:"certificates"`
	Valid           bool     `json:"is_valid"`
	ValidationError *string  `json:"validation_error"`
	CommonName      *string  `json:"common_name"`
	AltNames        []string `json:"alt_names"`
	Issuer          *string  `json:"issuer"`
}

// ServerKeyExchange represents the raw key data sent by the server in TLS key exchange message
type ServerKeyExchange struct {
	Key []byte `json:"key"`
}

// ServerFinished represents a TLS Finished message sent by the server
type ServerFinished struct {
	VerifyData []byte `json:"verify_data"`
}

// ServerHandshake stores all of the messages sent by the server during a standard TLS Handshake.
// It implements zgrab.EventData interface
type ServerHandshake struct {
	ServerHello        *ServerHello
	ServerCertificates *ServerCertificates
	ServerKeyExchange  *ServerKeyExchange
	ServerFinished     *ServerFinished
}

type encodedHandshake struct {
	Hello        *ServerHello        `json:"server_hello"`
	Certificates *ServerCertificates `json:"server_certificates"`
	KeyExchange  *ServerKeyExchange  `json:"server_key_exchange"`
	Finished     *ServerFinished     `json:"server_finished"`
}

// GetType always returns the TLS Handshake type
func (hs *ServerHandshake) GetType() EventType {
	return CONNECTION_EVENT_TLS
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
