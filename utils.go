package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

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

func printPacketInfo(packet gopacket.Packet) {
	// Let's see if the packet is an ethernet packet
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	log.Println("----printPacketInfo ...")
	log.Println(packet)
	if ethernetLayer != nil {
		fmt.Println("Ethernet layer detected.")
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
		fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
		// Ethernet type is typically IPv4 but could be ARP or other
		fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
		fmt.Println()
	}

	// Let's see if the packet is IP (even though the ether type told us)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)

		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		fmt.Println("Protocol: ", ip.Protocol)
		fmt.Println()
	}

	// Let's see if the packet is TCP
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)

		// TCP layer variables:
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		fmt.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
		fmt.Println("Sequence number: ", tcp.Seq)
		fmt.Println()
	}

	// Iterate over all layers, printing out each layer type
	fmt.Println("All packet layers:")
	for _, layer := range packet.Layers() {
		fmt.Println("- ", layer.LayerType())
	}

	// When iterating through packet.Layers() above,
	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		fmt.Println("Application layer/Payload found.")
		fmt.Printf("%s\n", applicationLayer.Payload())

		// Search for a string inside the payload
		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
			fmt.Println("HTTP found!")
		}
	}

	// Check for errors
	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}
}

func isClosed(c *net.TCPConn) bool {
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
