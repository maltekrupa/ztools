package zencoding

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func TestWithCheck(t *testing.T) { TestingT(t) }

type ZGrabSuite struct{}

var _ = Suite(&ZGrabSuite{})

func (s *ZGrabSuite) SetUpSuite(c *C) {
	RegisterEventType(mockEventType)
}

func (s *ZGrabSuite) TestUnknownEventTypeFromName(c *C) {
	_, err := EventTypeFromName("wat")
	c.Check(err, Not(IsNil))
}

func (s *ZGrabSuite) TestDecodeEmptyGrab(c *C) {
	g := new(Grab)
	g.Time = time.Unix(8675309, 0)
	g.Host = net.ParseIP("1.2.3.4")
	var d Grab
	marshalAndUnmarshal(g, &d, c)
}

func (s *ZGrabSuite) TestDecodeGrab(c *C) {
	g := new(Grab)
	g.Time = time.Unix(123456789, 0)
	g.Host = net.ParseIP("2.3.4.5")
	g.Log = make([]ConnectionEvent, 1)
	g.Log[0].Data = new(mockEventData).saneDefaults()
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
