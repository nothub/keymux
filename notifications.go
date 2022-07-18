package main

import (
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
	"log"
	"time"
)

type dbusIcon string

const (
	IconOn  dbusIcon = "media-playback-start"
	IconOff dbusIcon = "media-playback-stop"
)

func dbusConnect() (*dbus.Conn, error) {
	dbusConn, err := dbus.SessionBusPrivate()
	if err != nil {
		return nil, err
	}

	err = dbusConn.Auth(nil)
	if err != nil {
		return nil, err
	}

	err = dbusConn.Hello()
	if err != nil {
		return nil, err
	}

	return dbusConn, nil
}

func notifyInit() (func(message string, icon dbusIcon), func()) {
	dbusConn, err := dbusConnect()
	shutdownFunc := func() {
		err = dbusConn.Close()
		if err != nil {
			log.Println(err)
		}
	}
	if err != nil {
		log.Println(err)
		return func(message string, icon dbusIcon) {
			log.Println(message)
		}, shutdownFunc
	}
	return func(message string, icon dbusIcon) {
		_, err := notify.SendNotification(dbusConn, notify.Notification{
			Summary:       "PoeQol",
			Body:          message,
			AppIcon:       string(icon),
			ExpireTimeout: 2 * time.Second,
		})
		if err != nil {
			log.Println(err)
		}
	}, shutdownFunc
}
