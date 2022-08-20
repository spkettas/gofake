package common

import (
	"encoding/binary"
	"fmt"
	"net"
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
	Body    []byte
)

func init() {
	Body = []byte(fmt.Sprintf(format, len(content), content))
}

// Ip2Byte ip转为int
func Ip2Byte(sIp string) [4]byte {
	ip := net.ParseIP(sIp)
	if ip == nil {
		return [4]byte{0, 0, 0, 0}
	}

	ip = ip.To4()
	return [4]byte{ip[0], ip[1], ip[2], ip[3]}
}

// Ip2Int ip转为int
func Ip2Int(sIp string) uint32 {
	ip := net.ParseIP(sIp)
	if ip == nil {
		return 0
	}

	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

// Int2Ip int转化为ip
func Int2Ip(ipLong uint32) string {
	ipByte := make([]byte, 4)

	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}
