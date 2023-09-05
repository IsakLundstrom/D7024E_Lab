package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const BootstrapIp string = "10.10.0.2"
const BootstrapId = "0000000000000000000000000000000000000000"
const ipPrefix string = "10.10.0"

func IsBootstrap() bool {
	myIp, err := GetMyIp()
	if err != nil {
		log.Fatal(err)
	}
	return myIp == BootstrapIp
}

func GetMyIp() (string, error) {
	containerHostname, _ := os.Hostname()
	ips, _ := net.LookupIP(containerHostname)

	for _, ip := range ips {
		fmt.Println(ip.String())
		if strings.HasPrefix(ip.String(), ipPrefix) {
			return ip.String(), nil
		}
	}
	return "", errors.New("no ip")
}
