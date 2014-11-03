package zencoding

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"ztools/ztls"
)

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

// SetVersion sets the version and range checks for validity
func (sh *ServerHello) SetVersion(vers uint16) *ServerHello {
	if vers < ztls.VersionSSL30 || vers > ztls.VersionTLS12 {
		log.Panic("Invalid TLS version %d", vers)
	}
	sh.Version = vers
	return sh
}

// SetOCSP sets OCSP support to true
func (sh *ServerHello) SetOCSP() *ServerHello {
	sh.OcspStapling = true
	return sh
}

// SetHeartbeat sets Heartbeat support to true
func (sh *ServerHello) SetHeartbeat() *ServerHello {
	sh.HeartbeatSupported = true
	return sh
}

// PopulateRandom creates 32-byte TLS random field
func (sh *ServerHello) PopulateRandom() *ServerHello {
	sh.Random = make([]byte, 32)
	io.ReadFull(rand.Reader, sh.Random)
	return sh
}

type ServerCertificates struct {
	Certificates    [][]byte `json:"certificates"`
	Valid           bool     `json:"is_valid"`
	ValidationError *string  `json:"validation_error"`
	CommonName      *string  `json:"common_name"`
	AltNames        []string `json:"alt_names"`
	Issuer          *string  `json:"issuer"`
}

type ServerKeyExchange struct {
	Key []byte `json:"key"`
}

type ServerFinished struct {
	VerifyData []byte `json:"verify_data"`
}

type Handshake struct {
	ServerHello        *ServerHello        `json:"server_hello"`
	ServerCertificates *ServerCertificates `json:"server_certificates"`
	ServerKeyExchange  *ServerKeyExchange  `json:"server_key_exchange"`
	ServerFinished     *ServerFinished     `json:"server_finished"`
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
	return c
}

func decodeKeyExchange(raw map[string]interface{}) *ServerKeyExchange {
	kx := new(ServerKeyExchange)
	return kx
}

func decodeFinished(raw map[string]interface{}) *ServerFinished {
	f := new(ServerFinished)
	return f
}

func decodeHandshake(raw map[string]interface{}) *Handshake {
	h := new(Handshake)
	rawHello, helloPresent := raw["server_hello"]
	if helloPresent {
		hello, _ := rawHello.(map[string]interface{})
		h.ServerHello = decodeHello(hello)
	}
	rawCerts, certsPresent := raw["server_certificates"]
	if certsPresent {
		certs, _ := rawCerts.(map[string]interface{})
		h.ServerCertificates = decodeCertificates(certs)
	}
	rawSkx, skxPresent := raw["server_key_exchange"]
	if skxPresent {
		skx, _ := rawSkx.(map[string]interface{})
		h.ServerKeyExchange = decodeKeyExchange(skx)
	}
	rawFinished, finishedPresent := raw["server_finished"]
	if finishedPresent {
		finished, _ := rawFinished.(map[string]interface{})
		h.ServerFinished = decodeFinished(finished)
	}
	return h
}
