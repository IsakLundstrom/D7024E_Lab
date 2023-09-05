package kademlia

import (
	"fmt"
	"net"
	"os"
)

const bootstrapIp string = "10.10.0.2"

func IsBootstrap() bool {
	containerHostname, _ := os.Hostname()
	ips, _ := net.LookupIP(containerHostname)

	for _, ip := range ips {
		fmt.Println(ip.String())
		if ip.String() == bootstrapIp {
			return true
		}
	}
	return false
}
