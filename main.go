package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"time"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyUSB0")
	neopixel := firmata.NewNeopixelDriver(firmataAdaptor, "6", 1) // PD6
	led := gpio.NewLedDriver(firmataAdaptor, "12")                // PB4
	button := gpio.NewButtonDriver(firmataAdaptor, "7")           // PD7

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
	if err := neopixel.SetPixel(0, 0xFF0000); err != nil {
		fmt.Printf("neopixel set pixel: %v\n", err)
	}
	if err := neopixel.Show(); err != nil {
		fmt.Printf("neopixel show: %v\n", err)
	}

	for {
		fmt.Print(".")
		if err := led.On(); err != nil {
			fmt.Printf("led on: %v\n", err)
		}
		if err := neopixel.SetPixel(0, 0x00FF00); err != nil {
			fmt.Printf("neopixel set pixel: %v\n", err)
		}
		if err := neopixel.Show(); err != nil {
			fmt.Printf("neopixel show: %v\n", err)
		}
		time.Sleep(time.Second)

		if err := led.Off(); err != nil {
			fmt.Printf("led off: %v\n", err)
		}
		if err := neopixel.SetPixel(0, 0x000000); err != nil {
			fmt.Printf("neopixel set pixel: %v\n", err)
		}
		if err := neopixel.Show(); err != nil {
			fmt.Printf("neopixel show: %v\n", err)
		}
		time.Sleep(time.Second)
	}
}
