package main

import (
	"flag"
	"fmt"
	"github.com/MarinX/keylogger"
	"github.com/jezek/xgb"
	"github.com/nothub/keymux/buildinfo"
	"golang.org/x/exp/slices"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	delay, device, keys, window := parseFlags()

	{
		con, err := xgb.NewConn()
		if err != nil {
			log.Fatal(err)
		}
		x = con
		defer func() {
			x.Close()
		}()
	}

	kbd := kbdAttach(device)
	defer func() {
		err := kbd.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	events := kbd.Read()
	mutex := sync.Mutex{}
	for event := range events {

		// filter for keyboard events
		if event.Type != keylogger.EvKey {
			continue
		}

		// filter for key pressed down
		if !event.KeyPress() {
			continue
		}

		// filter flask hotkeys
		// TODO: else if pause hotkey
		if !slices.Contains(keys, event.KeyString()) {
			continue
		}

		// check active window
		if window != "" && window != activeWindow() {
			continue
		}

		// begin key multiplex
		if !mutex.TryLock() {
			continue
		}

		log.Printf("key pressed: %v", event.KeyString())
		for _, key := range keys {
			if key == event.KeyString() {
				continue
			}
			err := kbd.WriteOnce(key)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(delay())
		}

		// ignore events until delay awaited
		go func() {
			time.Sleep(delay())
			mutex.Unlock()
		}()
	}
}

func parseFlags() (func() time.Duration, string, []string, string) {
	flag.Usage = printUsage

	keysFlag := flag.String("keys", "1,2,3,4,5", "")
	deviceFlag := flag.String("device", "", "")
	windowFlag := flag.String("window", "", "")
	delayFlag := flag.Int("delay", 50, "")
	delayRandFlag := flag.Int("delay-random", 25, "")

	flag.Parse()

	keys := strings.Split(*keysFlag, ",")
	if len(keys) < 2 {
		log.Printf("Invalid keys argument: %q\n", keys)
		flag.Usage()
	}

	delay := time.Duration(*delayFlag) * time.Millisecond
	if delay < 10*time.Millisecond {
		log.Printf("Invalid delay argument: %q\n", delay)
		flag.Usage()
	}

	delayRand := time.Duration(*delayRandFlag) * time.Millisecond

	delayFunc := func() time.Duration {
		if delayRand.Milliseconds() <= 0 {
			return delay
		}
		n := delay.Nanoseconds()
		n = n + random.Int63n(delayRand.Nanoseconds())
		return time.Duration(n)
	}

	return delayFunc, *deviceFlag, keys, *windowFlag
}

func printUsage() {
	fmt.Print("keymux " + buildinfo.Version + "\n" +
		"\n" +
		"Multiplex multiple keyboard inputs with a single keypress.\n" +
		"Given a set of keys, if any key is pressed, keymux will simulate the presses of the other keys.\n" +
		"\n" +
		"Usage:\n" +
		"  keymux --keys \"a,b,x\"\n" +
		"  keymux --device \"/dev/input/by-id/usb-042-event-kbd\"\n" +
		"  keymux --window \"Path of Exile\"\n" +
		"  keymux --delay 50 --delay-random 30\n" +
		"\n" +
		"Options:\n" +
		"  --keys=<keys>\n" +
		"      List of hotkeys to monitor and send. (default: \"1,2,3,4,5\")\n" +
		"  --device=<path>\n" +
		"      Keyboard device path. (omit for auto search)\n" +
		"  --window=<name>\n" +
		"      Restrict input to specific window. (omit to disable)\n" +
		"  --delay=<ms>\n" +
		"      Keysend delay in milliseconds. (default: 50)\n" +
		"  --delay-random=<ms>\n" +
		"      Additional random keysend delay. (default: 25)\n" +
		"  -h --help\n" +
		"      Print help message and exit.\n")
}
