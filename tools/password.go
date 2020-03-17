package main

import (
	"encoding/hex"
	"flag"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Options struct {
	Password string
}

func main() {
	options := &Options{}

	flag.StringVar(&options.Password, "password", "", "password")

	flag.Parse()

	hashed, err := generateHashFromPassword(options.Password)
	if err != nil {
		panic(err)
	}

	log.Printf("%v %v", options.Password, hex.EncodeToString(hashed))
}

func generateHashFromPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
