package main

import (
	"fmt"

	"github.com/bahner/go-ma/key/set"
	nanoid "github.com/matoous/go-nanoid/v2"
	log "github.com/sirupsen/logrus"
)

func generateKeyset() string {

	name, err := nanoid.New()
	if err != nil {
		log.Fatalf("Failed to generate new nanoid: %v", err)
	}

	ks, err := set.New(name)
	if err != nil {
		log.Fatalf("Failed to generate new keyset: %v", err)
	}

	pks, err := ks.Pack()
	if err != nil {
		log.Fatalf("Failed to pack keyset: %v", err)
	}

	fmt.Println("export " + keysetEnvVar + "=" + pks)

	return pks
}
