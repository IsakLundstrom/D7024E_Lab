package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"net"
	"os"
	"strings"
)

const PROTOCOL string = "tcp"
const PORT string = "4000"
const IP_PREFIX string = "10.10.0"
const BOOTSTRAP_IP string = "10.10.0.2"
const BOOTSTRAP_ID string = "0000000000000000000000000000000000000000"

func IsBootstrap(prefix string) (bool, error) {
	myIp, err := GetMyIp(prefix)
	if err != nil {
		return false, err
	}
	return myIp == BOOTSTRAP_IP, nil
}

func GetMyIp(prefix string) (string, error) {
	containerHostname, _ := os.Hostname()
	ips, _ := net.LookupIP(containerHostname)

	for _, ip := range ips {
		if strings.HasPrefix(ip.String(), prefix) {
			return ip.String(), nil
		}
	}
	return "", errors.New("no ip")
}

func GetHash(data []byte) string {
	hasher := sha1.New()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
