package logger

import (
	"fmt"
	"log"
	"net"
	"os"
)

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

// Get hostname instance
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return hostname
}
