package zencoding

import (
	"encoding/json"
	"net"
	"time"

	. "gopkg.in/check.v1"
)

type ZGrabSuite struct{}

var _ = Suite(&ZGrabSuite{})

func (s *ZGrabSuite) TestEventTypeFromNameConnect(c *C) {
	connect, err := EventTypeFromName(CONNECTION_EVENT_CONNECT_NAME)
	c.Check(err, IsNil)
	c.Check(connect, Equals, CONNECTION_EVENT_CONNECT)
}

func (s *ZGrabSuite) TestEventTypeFromNameTLS(c *C) {
	tls, err := EventTypeFromName(CONNECTION_EVENT_TLS_NAME)
	c.Check(err, IsNil)
	c.Check(tls, Equals, CONNECTION_EVENT_TLS)
}

func (s *ZGrabSuite) TestUnknownEventTypeFromName(c *C) {
	_, err := EventTypeFromName("wat")
	c.Check(err, Not(IsNil))
}

func (s *ZGrabSuite) TestDecodeEventTypeConnect(c *C) {
	var t EventType
	orig := CONNECTION_EVENT_CONNECT
	marshalAndUnmarshal(&orig, &t, c)
}

func (s *ZGrabSuite) TestDecodeEventTypeTLS(c *C) {
	var t EventType
	orig := CONNECTION_EVENT_TLS
	marshalAndUnmarshal(&orig, &t, c)
}

func (s *ZGrabSuite) TestDecodeConnectionEvent(c *C) {
	h := new(ServerHandshake).saneDefaults()
	event := ConnectionEvent{
		Data:  h,
		Error: nil,
	}
	var decodedEvent ConnectionEvent
	marshalAndUnmarshal(&event, &decodedEvent, c)
}

func (s *ZGrabSuite) TestDecodeEmptyGrab(c *C) {
	g := new(Grab)
	g.Time = time.Unix(8675309, 0)
	var d Grab
	marshalAndUnmarshal(g, &d, c)
}

func (s *ZGrabSuite) TestDecodeGrab(c *C) {
	g := new(Grab)
	g.Time = time.Unix(123456578, 0)
	g.Domain = "davidadrian.org"
	event := ConnectionEvent{
		Error: nil,
		Data:  new(ServerHandshake).saneDefaults(),
	}
	g.Host = net.ParseIP("1.2.3.4")
	g.Log = []ConnectionEvent{event}
	var d Grab
	marshalAndUnmarshal(g, &d, c)
}

func marshalAndUnmarshal(original interface{}, target interface{}, c *C) {
	b, err := json.Marshal(original)
	c.Assert(err, IsNil)
	err = json.Unmarshal(b, target)
	c.Assert(err, IsNil)
	c.Check(target, DeepEquals, original)
}
