package ztls

type ServerHello struct {
	Version             uint16 `json:"version"`
	Random              []byte `json:"random"`
	SessionID           []byte `json:"session_id"`
	CipherSuite         uint16 `json:"cipher_suite"`
	CompressionMethod   uint8  `json:"compression_method"`
	OcspStapling        bool   `json:"ocsp_stapling"`
	TicketSupported     bool   `json:"ticket"`
	SecureRenogotiation bool   `json:"secure_renogotiation"`
	HeartbeatSupported  bool   `json:"heartbeat"`
}

// ServerCertificates represents a TLS certificates message in a format friendly to the golang JSON library.
// ValidationError should be non-nil whenever Valid is false.
type Certificates struct {
	Certificates    [][]byte
	Valid           bool
	ValidationError error
	CommonName      string
	AltNames        []string
	Issuer          string
}

// ServerKeyExchange represents the raw key data sent by the server in TLS key exchange message
type ServerKeyExchange struct {
	Key []byte `json:"key"`
}

// Finished represents a TLS Finished message
type Finished struct {
	VerifyData []byte `json:"verify_data"`
}

// ServerHandshake stores all of the messages sent by the server during a standard TLS Handshake.
// It implements zgrab.EventData interface
type ServerHandshake struct {
	ServerHello        *ServerHello       `json:"server_hello"`
	ServerCertificates *Certificates      `json:"server_certificates"`
	ServerKeyExchange  *ServerKeyExchange `json:"server_key_exchange"`
	ServerFinished     *Finished          `json:"server_finished"`
}

func (c *Conn) GetHandshakeLog() *ServerHandshake {
	return c.handshakeLog
}

func (m *serverHelloMsg) MakeLog() *ServerHello {
	sh := new(ServerHello)
	sh.Version = m.vers
	sh.Random = make([]byte, len(m.random))
	copy(sh.Random, m.random)
	sh.SessionID = make([]byte, len(m.sessionId))
	copy(sh.SessionID, m.sessionId)
	sh.CipherSuite = m.cipherSuite
	sh.CompressionMethod = m.compressionMethod
	sh.OcspStapling = m.ocspStapling
	sh.TicketSupported = m.ticketSupported
	sh.SecureRenogotiation = m.secureRenegotiation
	sh.HeartbeatSupported = m.heartbeat
	return sh
}

func (m *certificateMsg) MakeLog() *Certificates {
	sc := new(Certificates)
	sc.Certificates = make([][]byte, len(m.certificates))
	for idx, cert := range m.certificates {
		sc.Certificates[idx] = make([]byte, len(cert))
		copy(sc.Certificates[idx], cert)
	}
	return sc
}

func (m *serverKeyExchangeMsg) MakeLog() *ServerKeyExchange {
	skx := new(ServerKeyExchange)
	skx.Key = make([]byte, len(m.key))
	copy(skx.Key, m.key)
	return skx
}

func (m *finishedMsg) MakeLog() *Finished {
	sf := new(Finished)
	sf.VerifyData = make([]byte, len(m.verifyData))
	copy(sf.VerifyData, m.verifyData)
	return sf
}
