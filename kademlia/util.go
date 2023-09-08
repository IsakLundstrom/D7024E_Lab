package main

import (
	"crypto/sha1"
	"errors"
	"log"
	"net"
	"os"
	"strings"
)

const PROTOCOL string = "tcp"
const PORT string = "4000"
const IP_PREFIX string = "10.10.0"
const BOOTSTRAP_IP string = "10.10.0.2"
const BOOTSTRAP_ID string = "0000000000000000000000000000000000000000"

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
		if strings.HasPrefix(ip.String(), IP_PREFIX) {
			return ip.String(), nil
		}
	}
	return "", errors.New("no ip")
}

func GetHash(data []byte) KademliaID{
	hasher := sha1.New()
	hasher.Write([] byte(data))
	hash := hasher.Sum(nil)
	return *NewKademliaIDByte(hash)
}
