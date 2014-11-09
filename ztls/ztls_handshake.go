package ztls

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

func (c *Conn) GetHandshakeLog() *ServerHandshake {
	return c.handshakeLog
}
