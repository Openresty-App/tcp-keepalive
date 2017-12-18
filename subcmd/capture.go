package subcmd

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
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

func newStatisticData(p gopacket.Packet) CaptureTcp {
	d := &CaptureTcpData{}

	// Check for errors
	if err := p.ErrorLayer(); err != nil {
		log.Fatal("Error decoding some part of the packet:", err)
		return nil
	}
	d.T = p.Metadata().Timestamp

	ipLayer := p.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		d.SIP = ip.SrcIP
		d.DIP = ip.DstIP
	}

	// Let's see if the packet is TCP
	tcpLayer := p.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		d.Sport = uint16(tcp.SrcPort)
		d.Dport = uint16(tcp.DstPort)
		d.Seq = tcp.Seq
		d.Ack = tcp.Ack
	}

	return d
}

func tcpdump(device string, filter string, d chan Capture) {
	log.Printf("start tcpdmp device:%s, filter:%s", device, filter)
	snapshotLen := 65535
	promiscuous := false
	timeout := 30 * time.Second

	h, err := pcap.OpenLive(device, int32(snapshotLen), promiscuous, timeout)
	if err != nil {
		log.Printf("tcpdump: %s\n", err)
		return
	}
	defer h.Close()

	if filter != "" {
		if err := h.SetBPFFilter(filter); err != nil {
			log.Printf("tcpdump: %s\n", err)
		}
	}

	packetSource := gopacket.NewPacketSource(h, h.LinkType())
	for {
		packet, err := packetSource.NextPacket()
		if err == io.EOF {
			log.Println("Error: io.EOF")
			break
		} else if err != nil {
			continue
		}

		d <- newStatisticData(packet)
	}
}
