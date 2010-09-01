package manager

import (
	"testing"
	"log"
)

func TestProfile(t *testing.T) {
	ip, err := GetBroadcastAddr()
	if err != nil {
		log.Stderr(err)
		t.Fail()
	}
	log.Stdout(ip)
}
