package main

import (
	"crypto/rand"
	"io/ioutil"
)

func main() {

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("encryptionKey", key, 0644)
	if err != nil {
		panic(err)
	}
}
