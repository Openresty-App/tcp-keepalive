package main

import (
	"fmt"
	"net"
	"time"
)

type Capture interface{}

type CaptureIP interface {
	GetSIP() net.IP
	GetDIP() net.IP
	GetTime() time.Time
}

type CaptureTcp interface {
	CaptureIP
	GetSport() uint16
	GetDport() uint16
	GetSeq() uint32
	GetAck() uint32
}

type CaptureTcpData struct {
	SIP   net.IP
	DIP   net.IP
	Sport uint16
	Dport uint16
	Seq   uint32
	Ack   uint32
	T     time.Time
}

func (c *CaptureTcpData) GetSIP() net.IP {
	return c.SIP
}

func (c *CaptureTcpData) GetDIP() net.IP {
	return c.DIP
}

func (c *CaptureTcpData) GetSport() uint16 {
	return c.Sport
}

func (c *CaptureTcpData) GetDport() uint16 {
	return c.Dport
}

func (c *CaptureTcpData) GetTime() time.Time {
	return c.T
}

func (c *CaptureTcpData) GetSeq() uint32 {
	return c.Seq
}

func (c *CaptureTcpData) GetAck() uint32 {
	return c.Ack
}

func (c *CaptureTcpData) ToString() string {
	return fmt.Sprintf("%s %s:%d-->%s:%d, req:%d, ack:%d\n",
		c.T.String(), c.SIP, c.Sport,
		c.DIP, c.Dport,
		c.Seq, c.Ack)
}
