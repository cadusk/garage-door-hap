package main

import (
	"log"

	"github.com/brutella/hc"
	"github.com/brutella/hc/characteristic"
	"github.com/kkyr/fig"
	"periph.io/x/conn/v3/gpio"
)

var (
	device *Device
	acc    *GarageDoor
	config Config
)

func main() {
	if err := fig.Load(&config); err != nil {
		panic(err)
	}

	device = NewDevice(config.Device.Address)
	acc = NewGarageDoor(config.Accessory.SerialNumber, AppVersion)

	determineStatus()
	acc.TargetDoorState.OnValueRemoteUpdate(commandToTransition)
	registerSensorForMonitoring(device.sensorOpen)
	registerSensorForMonitoring(device.sensorClosed)

	startHomeKitServer()
}

func registerSensorForMonitoring(sensor gpio.PinIn) {
	go func() {
		for {
			sensor.WaitForEdge(-1)
			determineStatus()
		}
	}()
}

func commandToTransition(value int) {
	// transition to state
	device.Click(config.Device.Channel)
	if value == characteristic.TargetDoorStateOpen {
		acc.CurrentDoorState.UpdateValue(characteristic.CurrentDoorStateOpening)
	} else {
		acc.CurrentDoorState.UpdateValue(characteristic.CurrentDoorStateClosing)
	}

	// TODO: Check if we ever get to desired state.
	// TODO: Do something if times out and we're not there yet
}

func determineStatus() {
	isOpen := device.sensorOpenStatus()
	isClosed := device.sensorClosedStatus()

	if isOpen && !isClosed {
		acc.CurrentDoorState.UpdateValue(characteristic.CurrentDoorStateOpen)
		acc.TargetDoorState.UpdateValue(characteristic.TargetDoorStateOpen)

	} else if !isOpen && isClosed {
		acc.CurrentDoorState.UpdateValue(characteristic.CurrentDoorStateClosed)
		acc.TargetDoorState.UpdateValue(characteristic.TargetDoorStateClosed)

	} else if !isOpen && !isClosed {
		currentDoorState := acc.CurrentDoorState.GetValue()
		if currentDoorState != characteristic.CurrentDoorStateOpen &&
			currentDoorState != characteristic.CurrentDoorStateClosed {
			// if we're in the middle of something, let it finish.
			return
		}

		targetDoorState := acc.TargetDoorState.GetValue()
		if targetDoorState == characteristic.TargetDoorStateClosed {
			acc.TargetDoorState.UpdateValue(characteristic.TargetDoorStateOpen)
		} else {
			acc.TargetDoorState.UpdateValue(characteristic.TargetDoorStateClosed)
		}

	} else {
		// TODO: This means that both sensors are connected
		// FIXME: Figure out what to do in this situation!!
	}
}

func startHomeKitServer() {
	hcConfig := hc.Config{
		Pin:         config.Accessory.Pin,
		StoragePath: config.Accessory.StoragePath,
		Port:        config.Accessory.Port,
	}

	t, err := hc.NewIPTransport(hcConfig, acc.Accessory)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		log.Println("Terminating...")
		device.Halt()

		<-t.Stop()
	})

	t.Start()
}
