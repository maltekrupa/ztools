package zencoding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"ztools/ztls"
)

func (sh *ServerHello) saneDefaults() *ServerHello {
	sh.SetVersion(ztls.VersionTLS12)
	sh.PopulateRandom()
	sh.SessionID = nil
	sh.CipherSuite = ztls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
	sh.CompressionMethod = 0
	sh.OcspStapling = false
	sh.TicketSupported = false
	sh.HeartbeatSupported = false
	return sh
}

// compare compares two ServerHellos for equality
func (a *ServerHello) compare(b *ServerHello) error {
	if a.Version != b.Version {
		return fmt.Errorf("Mismatched Versions. Expected %d, had %d", a.Version, b.Version)
	}
	if !bytes.Equal(a.Random, b.Random) {
		return fmt.Errorf("Mismatched Randoms.")
	}
	if !bytes.Equal(a.SessionID, b.SessionID) {
		return fmt.Errorf("Mismatched Session ID.")
	}
	if a.CipherSuite != b.CipherSuite {
		return fmt.Errorf("Mismatched Cipher Suites. Expected %d, had %d", a.CipherSuite, b.CipherSuite)
	}
	if a.CompressionMethod != b.CompressionMethod {
		return fmt.Errorf("Mismatched Compression Method. Expected %d, had %d", a.CompressionMethod, b.CompressionMethod)
	}
	if a.OcspStapling != b.OcspStapling {
		return fmt.Errorf("Mismatched OCSP. Expected %b, had %b", a.OcspStapling, b.OcspStapling)
	}
	if a.TicketSupported != b.TicketSupported {
		return fmt.Errorf("Mismatched Ticket. Expected %b, had %b", a.TicketSupported, b.TicketSupported)
	}
	if a.HeartbeatSupported != b.HeartbeatSupported {
		return fmt.Errorf("Mismatched Heartbeat. Expected %b, had %b", a.HeartbeatSupported, b.HeartbeatSupported)
	}
	return nil
}

func testHelloHelper(sh *ServerHello, t *testing.T) {
	// Serialize it
	serialized, err := json.Marshal(sh)
	if err != nil {
		t.Error(err)
	}
	// Deserialize it
	var r interface{}
	json.Unmarshal(serialized, &r)
	d := decodeHello(r.(map[string]interface{}))
	// Compare
	if cmpErr := sh.compare(d); cmpErr != nil {
		t.Log(cmpErr)
		t.Fail()
	}
}

func TestDecodeHello(t *testing.T) {
	sh := new(ServerHello)
	sh.saneDefaults()
	testHelloHelper(sh, t)
}

func TestDecodeHelloComplicated(t *testing.T) {
	sh := new(ServerHello)
	sh.saneDefaults().SetVersion(ztls.VersionSSL30).SetOCSP().SetHeartbeat()
	testHelloHelper(sh, t)
}
