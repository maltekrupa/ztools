package zlog

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestLogger(t *testing.T) { TestingT(t) }

type LoggerSuite struct{}

var _ = Suite(&LoggerSuite{})

func (s *LoggerSuite) TestPrint(c *C) {
	Print(LOG_DEBUG, 2014, "yay", "research")
}

func (s *LoggerSuite) TestPrintf(c *C) {
	Printf(LOG_TRACE, "It's %d, %s %s %s %s?", 2019, "are", "we", "done", "yet")
}
