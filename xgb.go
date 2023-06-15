package main

import (
	"fmt"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"log"
	"strings"
)

var x *xgb.Conn

func activeWindow() string {
	// https://github.com/BurntSushi/xgb/blob/deaf085860bc2ec870e50cacc83c467b3501a404/examples/get-active-window/main.go

	// Get the window id of the root window.
	setup := xproto.Setup(x)
	root := setup.DefaultScreen(x).Root

	// Get the atom id (i.e., intern an atom) of "_NET_ACTIVE_WINDOW".
	aname := "_NET_ACTIVE_WINDOW"
	activeAtom, err := xproto.InternAtom(x, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		log.Fatal(err)
	}

	// Get the atom id (i.e., intern an atom) of "_NET_WM_NAME".
	aname = "_NET_WM_NAME"
	nameAtom, err := xproto.InternAtom(x, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		log.Fatal(err)
	}

	// Get the actual value of _NET_ACTIVE_WINDOW.
	// Note that 'reply.Value' is just a slice of bytes, so we use an
	// XGB helper function, 'Get32', to pull an unsigned 32-bit integer out
	// of the byte slice. We then convert it to an X resource id so it can
	// be used to get the name of the window in the next GetProperty request.
	reply, err := xproto.GetProperty(x, false, root, activeAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		log.Fatal(err)
	}
	windowId := xproto.Window(xgb.Get32(reply.Value))
	fmt.Printf("Active window id: %X\n", windowId)

	// Now get the value of _NET_WM_NAME for the active window.
	// Note that this time, we simply convert the resulting byte slice,
	// reply.Value, to a string.
	reply, err = xproto.GetProperty(x, false, windowId, nameAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Trim(string(reply.Value), " \n")
}
