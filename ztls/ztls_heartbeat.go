package ztls

import (
	"errors"
)

const (
	// Record Type
	recordTypeHeartbeat recordType = 24

	// Extension Number
	extensionHeartbeat uint16 = 15

	// Heartbeat Mode
	heartbeatModePeerAllowed    uint8 = 1
	heartbeatModePeerNotAllowed uint8 = 2

	// Heartbeat Message Types
	heartbeatTypeRequest  uint8 = 1
	heartbeatTypeResponse uint8 = 2
)

var (
	HeartbleedError = errors.New("Error after Heartbleed")
)

type Heartbleed struct {
	HeartbeatEnabled bool
	Vulnerable       bool
}

type heartbleedMessage struct {
	raw []byte
}

func (m *heartbleedMessage) marshal() []byte {
	x := make([]byte, 3)
	x[0] = 1
	x[1] = byte(0x00)
	x[2] = byte(0x00)
	m.raw = x
	return x
}

func (c *Conn) CheckHeartbleed(b []byte) (n int, err error) {
	if err = c.Handshake(); err != nil {
		return
	}
	if !c.heartbeat {
		return
	}
	c.in.Lock()
	defer c.in.Unlock()

	hb := heartbleedMessage{}
	hb.marshal()

	if _, err = c.writeRecord(recordTypeHeartbeat, hb.raw); err != nil {
		return 0, err
	}

	if err = c.readRecord(recordTypeHeartbeat); err != nil {
		return 0, HeartbleedError
	}
	if err = c.in.error(); err != nil {
		return 0, HeartbleedError
	}
	n, err = c.input.Read(b)
	if c.input.off >= len(c.input.data) {
		c.in.freeBlock(c.input)
		c.input = nil
	}

	if n != 0 {
		return n, HeartbleedError
	}
	if err != nil {
		return 0, HeartbleedError
	}
	return 0, HeartbleedError
}