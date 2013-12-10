package main

import (
	"flag"
	"fmt"
	"github.com/guelfey/go.dbus"
	"os"
	"os/exec"
	"time"
)

type Options struct {
	time_action      float64
	polling_interval int
	repeat           bool
	action           string
}

type UPower struct {
	sysBusConn *dbus.Conn
	battery    *dbus.Object
}

func (self *UPower) connect() {
	var e error
	self.sysBusConn, e = dbus.SystemBus()
	if e != nil {
		panic(e)
	}
	self.battery = self.sysBusConn.Object("org.freedesktop.UPower", "/org/freedesktop/UPower/devices/battery_BAT0")
}

func (self *UPower) getTimeToEmpty() (timeToEmpty float64) {
	v, e := self.battery.GetProperty("org.freedesktop.UPower.Device.TimeToEmpty")
	if e != nil {
		panic(e)
	}
	switch v.Value().(type) {
	case int32:
		timeToEmpty = float64(v.Value().(int32)) / 60
	case int64:
		timeToEmpty = float64(v.Value().(int64)) / 60
	}
	return
}

func (self *Options) parseCommandFlags() {
	flag.Float64Var(&self.time_action, "time-action", 10, "The time remaining in minutes of the battery when critical action is taken.")
	flag.IntVar(&self.polling_interval, "polling-interval", 10, "The time remaining in minutes of the battery when critical action is taken.")
	flag.BoolVar(&self.repeat, "repeat", false, "The time remaining in minutes of the battery when critical action is taken.")
	flag.StringVar(&self.action, "action", "echo 'Situation Critical!'", "The time remaining in minutes of the battery when critical action is taken.")

	flag.Float64Var(&self.time_action, "t", 10, "The time remaining in minutes of the battery when critical action is taken. (shorthand)")
	flag.IntVar(&self.polling_interval, "p", 10, "The time remaining in minutes of the battery when critical action is taken. (shorthand)")
	flag.BoolVar(&self.repeat, "r", false, "The time remaining in minutes of the battery when critical action is taken. (shorthand)")
	flag.StringVar(&self.action, "a", "echo 'Situation Critical!'", "The time remaining in minutes of the battery when critical action is taken. (shorthand)")

	flag.Parse()
}

func main() {
	// parse command flags before doing anything else, in case of -h
	var options Options
	options.parseCommandFlags()
	fmt.Println(options)

	var battery UPower
	var timeToEmpty float64
	var cmd *exec.Cmd
	var execute = true

	battery.connect()

	for {
		timeToEmpty = battery.getTimeToEmpty()

		if timeToEmpty < options.time_action {
			if execute {
				cmd = exec.Command("sh", "-c", options.action)
				cmd.Stdout = os.Stdout
				cmd.Start()
				execute = options.repeat
			}
		} else {
			execute = true
		}

		time.Sleep(10 * time.Second)
	}
}
