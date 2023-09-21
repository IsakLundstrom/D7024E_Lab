package kademlia

import (
	"net"
	"os"
	"testing"
)

func TestIsBootstrap(t *testing.T) {
	containerHostname, _ := os.Hostname()
	ips, _ := net.LookupIP(containerHostname)

	expected := false
	isBoot, err := IsBootstrap(ips[0].String())
	if err != nil {
		t.Errorf("got an error: %s", err)
	}
	if expected != isBoot {
		t.Errorf("expected false but got true")
	}

	isBoot, err = IsBootstrap("not a prefix")
	if err == nil {
		t.Errorf("expected an error but didnt get one")
	}
	if expected != isBoot {
		t.Errorf("expected false but got true")
	}
}

func TestGetMyIp(t *testing.T) {
	containerHostname, _ := os.Hostname()
	ips, _ := net.LookupIP(containerHostname)

	expectedIp := ips[0].String()
	testIp, _ := GetMyIp(ips[0].String())
	if expectedIp != testIp {
		t.Errorf("expected ip %s but got %s", expectedIp, testIp)
	}
}

func TestGetHash(t *testing.T) {
	testData := "Hello I am test data"
	expectedHash := "7704b4ddde293816eee2bb041e9ddd4b71611a3a"
	generatedHash := GetHash([]byte(testData))
	if generatedHash != expectedHash {
		t.Errorf("expected hash %s but got %s", expectedHash, generatedHash)
	}
}
