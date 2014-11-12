package zencoding

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestEventSuite(t *testing.T) { TestingT(t) }

type EventSuite struct{}

var _ = Suite(&EventSuite{})

func (s *EventSuite) TestEventTypeFromName(c *C) {
	_, err := EventTypeFromName("does-not-exist")
	c.Check(err, NotNil)
}
