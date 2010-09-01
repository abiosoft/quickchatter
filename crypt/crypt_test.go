package crypt

import (
	"testing"
)

func TestKey(t *testing.T) {
	key := GenerateKey(16)
	if len(key) != 16 {
		t.Fail()
	}
}

func TestCrypt(t *testing.T) {
	key := GenerateKey(16)
	b := GenerateKey(64)
	c := make([]byte, 64)
	copy(c, b)
	Encrypt(key, b)
	Decrypt(key, b)
	if string(b) != string(c) {
		t.Fail()
	}
}

func TestMd5(t *testing.T) {
	hash := "bc9f8081feeaf2389e3b918d3366114f"
	s := "hello over there"
	if hash != Md5([]byte(s)) {
		t.Fail()
	}
}
