package main

import (
	"fmt"
	"time"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

var neopixel *firmata.NeopixelDriver

func showLed(r, g, b int) error {
	color := ((r & 0xFF) << 16) + ((g & 0xFF) << 8) + ((b & 0xFF) << 0)
	if err := neopixel.SetPixel(0, uint32(color)); err != nil {
		return fmt.Errorf("neopixel set pixel: %v\n", err)
	}
	if err := neopixel.Show(); err != nil {
		return fmt.Errorf("neopixel show: %v\n", err)
	}
	return nil
}

func setup() {
	firmataAdaptor := firmata.NewAdaptor("/dev/tty.usbserial-FTGCVV4Y")
	neopixel = firmata.NewNeopixelDriver(firmataAdaptor, "6", 1) // PD6
	led := gpio.NewLedDriver(firmataAdaptor, "12")               // PB4
	button := gpio.NewButtonDriver(firmataAdaptor, "7")          // PD7

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{neopixel, led, button},
	)
	button.On(gpio.ButtonPush, func(data interface{}) {
		fmt.Println("Release")
		if err := led.On(); err != nil {
			fmt.Printf("led on: %v\n", err)
		}
	})
	button.On(gpio.ButtonRelease, func(data interface{}) {
		fmt.Println("Press")
		if err := led.Off(); err != nil {
			fmt.Printf("led off: %v\n", err)
		}
	})

	robot.Start(false)
	if err := showLed(0, 0, 0); err != nil {
		fmt.Printf("%v\n", err)
	}
}

const (
	minLedValue = 0x00
	maxLedValue = 0xFF
	offsetLedValue = 0x08
	wait = time.Millisecond * 50
	redMask = 0x01 << 0
	greenMask = 0x01 << 1
	blueMask = 0x01 << 2
)

func showLedUpWave(r, g, b, mask int) error {
	for v := minLedValue; v < maxLedValue; v += offsetLedValue {
		if ((mask & redMask) != 0) {
			r = v
		}
		if ((mask & greenMask) != 0) {
			g = v
		}
		if ((mask & blueMask) != 0) {
			b = v
		}
		if err := showLed(r, g, b); err != nil {
			return err
		}
		time.Sleep(wait)
	}
	return nil
}

func showLedDownWave(r, g, b, mask int) error {
	for v := minLedValue; v < maxLedValue; v += offsetLedValue {
		if ((mask & redMask) != 0) {
			r = maxLedValue - v
		}
		if ((mask & greenMask) != 0) {
			g = maxLedValue - v
		}
		if ((mask & blueMask) != 0) {
			b = maxLedValue - v
		}
		if err := showLed(r, g, b); err != nil {
			return err
		}
		time.Sleep(wait)
	}
	return nil
}

func loop() {
	showLedUpWave(minLedValue,   maxLedValue, maxLedValue, greenMask) // gb
	showLedDownWave(minLedValue, maxLedValue, minLedValue, blueMask)  // g
	showLedUpWave(maxLedValue,   maxLedValue, minLedValue, redMask)   // rg
	showLedDownWave(maxLedValue, maxLedValue, minLedValue, greenMask) // r
	showLedUpWave(maxLedValue,   minLedValue, minLedValue, blueMask)  // rb
	showLedDownWave(minLedValue, minLedValue, maxLedValue, redMask)   // b
}

func main() {
	setup()
	showLedUpWave(minLedValue,   minLedValue, maxLedValue, blueMask)  // b
	for {
		loop()
	}
}
