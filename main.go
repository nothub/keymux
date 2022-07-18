package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/MarinX/keylogger"
)

func main() {
	delay, device, keys := settings()

	// TODO: replace shutdown callbacks with channels
	keyboard, shutdownKeyboard := keyboardAttach(device)
	defer shutdownKeyboard()
	notify, shutdownNotify := notifyInit()
	defer shutdownNotify()

	notify("Flask Multiplexer", IconOn)

	events := keyboard.Read()
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
		if activeWindow() != "Path of Exile" {
			continue
		}
		// begin key multiplex
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
		// ignore events until delay awaited
		go func() {
			time.Sleep(delay)
			mutex.Unlock()
		}()
	}

	notify("Flask Multiplexer", IconOff)

}

func keyboardAttach(device string) (*keylogger.KeyLogger, func()) {
	log.Println("Listening to keyboard", device)
	keyboard, err := keylogger.New(device)
	if err != nil {
		log.Fatal(err)
	}
	shutdownFunc := func() {
		err := keyboard.Close()
		if err != nil {
			log.Println(err)
		}
	}
	return keyboard, shutdownFunc
}

func activeWindow() string {
	cmd := exec.Command("xdotool", "getwindowfocus", "getwindowname") // TODO: do this without xdotool
	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	err := cmd.Run()
	stdErrTrim := strings.Trim(stdErr.String(), " \n")
	if len(stdErrTrim) > 0 {
		log.Println("stderr:", stdErrTrim)
	}
	if err != nil {
		log.Println(err)
		return ""
	}
	return strings.Trim(stdOut.String(), " \n")
}

func settings() (time.Duration, string, []string) {
	flag.Usage = func() {
		infos, ok := debug.ReadBuildInfo()
		name := infos.Main.Path + " " + infos.Main.Version
		if !ok {
			name = os.Args[0]
		}
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s\n", name)
		flag.PrintDefaults()
		os.Exit(0)
	}

	delayFlag := flag.Int("delay", 50, "keysend delay in milliseconds (min value 10)")
	deviceFlag := flag.String("device", "/dev/input/by-id/usb-046a_0011-event-kbd", "keyboard device path")
	keysFlag := flag.String("keys", "1,2,3,4,5", "list of 5 flask hotkeys")

	flag.Parse()

	delay := time.Duration(*delayFlag) * time.Millisecond
	if delay < 10*time.Millisecond {
		log.Println("Invalid argument", "delay")
		flag.Usage()
	}

	device := *deviceFlag
	if len(device) == 0 {
		log.Println("Invalid argument", "device")
		flag.Usage()
	}

	keys := strings.Split(*keysFlag, ",")
	if len(keys) != 5 {
		log.Println("Invalid argument", "keys")
		flag.Usage()
	}

	return delay, device, keys
}
