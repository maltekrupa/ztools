package zencoding

import (
	"encoding/base64"
	"encoding/json"
)

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

// GetType always returns the TLS Handshake type
func (hs *ServerHandshake) GetType() EventType {
	return typeTLSInstance
}

// UnpackMap extracts a map[string]interface{} into a ServerHandshake struct
func (hs *ServerHandshake) UnpackMap(map[string]interface{}) error {
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (hs *ServerHandshake) MarshalJSON() ([]byte, error) {
	// Prevent infinite recursion
	obj := struct {
		Hello        *ServerHello        `json:"server_hello"`
		Certificates *ServerCertificates `json:"server_certificates"`
		KeyExchange  *ServerKeyExchange  `json:"server_key_exchange"`
		Finished     *ServerFinished     `json:"server_finished"`
	}{
		Hello:        hs.ServerHello,
		Certificates: hs.ServerCertificates,
		KeyExchange:  hs.ServerKeyExchange,
		Finished:     hs.ServerFinished,
	}
	return json.Marshal(obj)
}

func decodeHello(raw map[string]interface{}) *ServerHello {
	h := new(ServerHello)
	h.Version = uint16(raw["version"].(float64))
	random := raw["random"].(string)
	h.Random, _ = base64.StdEncoding.DecodeString(random)
	if raw["session_id"] != nil {
		sid := raw["session_id"].(string)
		h.SessionID, _ = base64.StdEncoding.DecodeString(sid)
	}
	h.CipherSuite = uint16(raw["cipher_suite"].(float64))
	h.CompressionMethod = uint8(raw["compression_method"].(float64))
	h.OcspStapling = raw["ocsp_stapling"].(bool)
	h.TicketSupported = raw["ticket_supported"].(bool)
	h.HeartbeatSupported = raw["heartbeat_supported"].(bool)
	return h
}

func decodeCertificates(raw map[string]interface{}) *ServerCertificates {
	c := new(ServerCertificates)
	if raw["certificates"] != nil {
		certs := raw["certificates"].([]interface{})
		c.Certificates = make([][]byte, len(certs))
		for idx, cert := range certs {
			c.Certificates[idx], _ = base64.StdEncoding.DecodeString(cert.(string))
		}
	}
	c.Valid = raw["is_valid"].(bool)
	c.ValidationError = getStringPointer(raw, "validation_error")
	c.CommonName = getStringPointer(raw, "common_name")
	c.AltNames = getStringArray(raw, "alt_names")
	c.Issuer = getStringPointer(raw, "issuer")
	return c
}

func decodeKeyExchange(raw map[string]interface{}) *ServerKeyExchange {
	skx := new(ServerKeyExchange)
	skx.Key = getBytes(raw, "key")
	return skx
}

func decodeFinished(raw map[string]interface{}) *ServerFinished {
	sf := new(ServerFinished)
	sf.VerifyData = getBytes(raw, "verify_data")
	return sf
}

func decodeServerHandshake(raw map[string]interface{}) *ServerHandshake {
	h := new(ServerHandshake)
	rawHello, helloPresent := raw["server_hello"]
	if helloPresent && rawHello != nil {
		hello, _ := rawHello.(map[string]interface{})
		h.ServerHello = decodeHello(hello)
	}
	rawCerts, certsPresent := raw["server_certificates"]
	if certsPresent && rawCerts != nil {
		certs, _ := rawCerts.(map[string]interface{})
		h.ServerCertificates = decodeCertificates(certs)
	}
	rawSkx, skxPresent := raw["server_key_exchange"]
	if skxPresent && rawSkx != nil {
		skx, _ := rawSkx.(map[string]interface{})
		h.ServerKeyExchange = decodeKeyExchange(skx)
	}
	rawFinished, finishedPresent := raw["server_finished"]
	if finishedPresent && rawFinished != nil {
		finished, _ := rawFinished.(map[string]interface{})
		h.ServerFinished = decodeFinished(finished)
	}
	return h
}

type typeTLS uint8

func (t typeTLS) String() string {
	return "tls_handshake"
}

const (
	typeTLSInstance typeTLS = 0
)
