package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	device                     = "eth1"
	snapshot_len int32         = 65535
	promiscuous  bool          = false
	timeout      time.Duration = 10 * time.Millisecond
	handle       *pcap.Handle
	err          error

	options gopacket.SerializeOptions
)

// capture libpcap抓包
func capture() {
	handle, err := pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	filter := "tcp and port 80"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Only capturing TCP port 80 packets")

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// 发送伪装包
		SendHijack(handle, packet)
	}
}

func main() {
	capture()
}
