package utils

import (
	"io"
	"net"
	"time"
)

func IsClosed(c *net.TCPConn) bool {
	one := make([]byte, 1)
	c.SetReadDeadline(time.Time{})
	_, err := c.Read(one)
	if err == io.EOF {
		c.Close()
		c = nil
		return true
	} else {
		var zero time.Time
		c.SetReadDeadline(zero)

		return false
	}
}
