package main

import (
	"github.com/MarinX/keylogger"
	"log"
	"path/filepath"
)

func kbdFind() string {
	pattern := "/dev/input/by-path/*-event-kbd"

	devices, err := filepath.Glob(pattern)
	if err != nil {
		log.Println(err)
	}

	if len(devices) < 1 {
		log.Fatalf("no keyboards found at path: %q\n", pattern)
	}

	if len(devices) > 1 {
		log.Printf("multiple keyboards found: %q\n", devices)
	}

	return devices[0]
}

func kbdAttach(device string) *keylogger.KeyLogger {
	if device == "" {
		device = kbdFind()
	}

	log.Printf("attaching to keyboard: %q\n", device)
	kbd, err := keylogger.New(device)
	if err != nil {
		log.Fatal(err)
	}

	return kbd
}
