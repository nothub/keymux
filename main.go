package main

import (
	"github.com/MarinX/keylogger"
	"golang.org/x/exp/slices"
	"log"
	"sync"
	"time"
)

func main() {

	delay := 100 * time.Millisecond
	device := "/dev/input/by-id/usb-046a_0011-event-kbd"
	keys := []string{"1", "2", "3", "4", "5"}

	keyboard, err := keylogger.New(device)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := keyboard.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	events := keyboard.Read()
	mutex := sync.Mutex{}
	for event := range events {
		if event.Type != keylogger.EvKey {
			continue
		}
		if !event.KeyPress() {
			continue
		}
		if slices.Contains(keys, event.KeyString()) {
			if !mutex.TryLock() {
				continue
			}
			log.Printf("flask key pressed: %v", event.KeyString())
			for _, key := range keys {
				if key == event.KeyString() {
					continue
				}
				err := keyboard.WriteOnce(key)
				if err != nil {
					log.Println(err)
				}
				time.Sleep(delay)
			}
			go func() {
				time.Sleep(delay)
				mutex.Unlock()
			}()
		}
	}

}
