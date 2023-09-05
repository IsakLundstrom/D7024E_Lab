package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const BOOTSTRAP_ID = "0000000000000000000000000000000000000000"
const BOOTSTRAP_IP string = "10.10.0.2"
const IP_PREFIX string = "10.10.0"

func IsBootstrap() bool {
	myIp, err := GetMyIp()
	if err != nil {
		log.Fatal(err)
	}
	return myIp == BOOTSTRAP_IP
}

func GetMyIp() (string, error) {
	containerHostname, _ := os.Hostname()
	ips, _ := net.LookupIP(containerHostname)

	for _, ip := range ips {
		fmt.Println(ip.String())
		if strings.HasPrefix(ip.String(), IP_PREFIX) {
			return ip.String(), nil
		}
	}
	return "", errors.New("no ip")
}
