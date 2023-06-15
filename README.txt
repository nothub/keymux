keymux v0.2.1

Multiplex multiple keyboard inputs with a single keypress.
Given a set of keys, if any key is pressed, keymux will simulate the presses of the other keys.

Usage:
  keymux --keys   "a,b,x"
  keymux --device "/dev/input/by-id/usb-042-event-kbd"
  keymux --window "Path of Exile"
  keymux --delay 70 --delay-random 30

Options:
  --keys=<keys>
      List of hotkeys to monitor and send. (default: "1,2,3,4,5")
  --key-pause=<key>
      Hotkey for pausing execution. (default: "Pause")
  --device=<path>
      Keyboard device path. (omit for auto search)
  --window=<name>
      Restrict input to specific window. (omit to disable)
  --delay=<ms>
      Keysend delay in milliseconds. (default: 50)
  --delay-random=<ms>
      Additional random keysend delay. (default: 25)
  -h --help
      Print help message and exit.
