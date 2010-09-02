package crypt

import (
	"rand"
	"time"
	"crypto/aes"
	"crypto/md5"
	"log"
	"fmt"
)

func init() {
	rand.Seed(time.Seconds())
}

func GenerateKey(l int) (key []byte) {
	key = make([]byte, l)
	for i := 0; i < l; i++ {
		key[i] = uint8(rand.Intn(256))
	}
	return
}

func GenerateNums(l int) (num string) {
	for i := 0; i < l; i++ {
		num += fmt.Sprint(int(uint8(rand.Intn(256))))
	}
	return num
}

func Encrypt(key, b []byte) {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Stderr(err.String())
		return
	}
	c.Encrypt(b, b)
}

func Decrypt(key, b []byte) {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Stderr(err.String())
		return
	}
	c.Decrypt(b, b)
}

func Md5(b []byte) string {
	hash := md5.New()
	hash.Write(b)
	return fmt.Sprintf("%x", hash.Sum())
}
