package utils

import (
	"log"

	entity "github.com/blanvam/rasp-garden/entities"
	rpio "github.com/stianeikeland/go-rpio"
)

// GPIOOutChange function to open or close Raspberry Pi GPIO output pin
func GPIOOutChange(pinNumber int, open bool) error {
	err := rpio.Open()
	defer rpio.Close()
	log.Println("GPIOOutChange")
	if err != nil {
		log.Println("ERROR")
		return entity.ErrRGPIO
	}
	pin := rpio.Pin(pinNumber)
	pin.Output() // pin.Mode(rpio.Output)
	if open == true {
		pin.High()
	} else {
		pin.Low()
	}
	return nil
}

// GPIOInChange function to open or close Raspberry Pi GPIO input pin
func GPIOInChange(pinNumber int, read bool) error {
	err := rpio.Open()
	defer rpio.Close()
	log.Println("GPIOInChange")
	if err != nil {
		log.Println("ERROR")
		return entity.ErrRGPIO
	}
	pin := rpio.Pin(pinNumber)
	pin.Input() // pin.Mode(rpio.Input)
	if read == true {
		pin.Read()
	} else {
		// TODO
	}
	return nil
}
