package main

import (
	"log"

	"github.com/brutella/hc"
	"github.com/brutella/hc/characteristic"
	"periph.io/x/conn/v3/gpio"
)

var (
	device *Device
	acc    *GarageDoor
	config Config
)

func main() {
	if err := loadConfig(&config); err != nil {
		panic(err)
	}

	device = NewDevice(config.Device.Address)
	acc = NewGarageDoor(config.Accessory.SerialNumber, AppVersion)

	acc.TargetDoorState.OnValueRemoteUpdate(targetTransitionCommand)
	registerSensorForMonitoring(device.sensorOpen)
	registerSensorForMonitoring(device.sensorClosed)

	determineStatus()

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

func targetTransitionCommand(newState int) {
	device.Click(config.Device.Channel)

	// indicates we're transitioning...
	if newState == characteristic.TargetDoorStateOpen {
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
		log.Println("determineStatus - isOpen")
		updateState(characteristic.CurrentDoorStateOpen)

	} else if !isOpen && isClosed {
		log.Println("determineStatus - isClosed")
		updateState(characteristic.CurrentDoorStateClosed)

	} else if !isOpen && !isClosed {
		log.Println("determineStatus - not fully closed or open")
		currentDoorState := acc.CurrentDoorState.GetValue()
		targetDoorState := acc.TargetDoorState.GetValue()

		log.Println("determineStatus - current", currentDoorState, "target", targetDoorState)

		switch currentDoorState {
		case characteristic.CurrentDoorStateOpening, characteristic.CurrentDoorStateClosing:
			// TODO: Figure out what to test here to ensure we're in the middle
			// of something and not just lost. add a timer maybe to validate for
			// how long we're in the same moving state before we decide to
			// notify the user.
			return
		}

		if targetDoorState == characteristic.TargetDoorStateClosed {
			log.Println("determineStatus - setting target state to open")
			acc.TargetDoorState.UpdateValue(characteristic.TargetDoorStateOpen)
		} else {
			log.Println("determineStatus - setting target state to closed")
			acc.TargetDoorState.UpdateValue(characteristic.TargetDoorStateClosed)
		}

	} else {
		log.Println("determineStatus - who knows?!?")
		// TODO: This means that both sensors are connected
		// FIXME: Figure out what to do in this situation!!
	}
}

func updateState(toState int) {
	oldState := acc.CurrentDoorState.GetValue()

	// TODO: should we worry about this scenario?
	// if oldState == toState {
	// 	log.Println("updateState - same state, leaving...")
	// 	return
	// }

	log.Println("updateState - switchting from ", oldState, "to", toState)
	switch toState {
	case characteristic.CurrentDoorStateClosing,
		characteristic.CurrentDoorStateOpening,
		characteristic.CurrentDoorStateStopped:

		log.Println("updateState - [ 2, 3, 4 ]:", toState)
		acc.CurrentDoorState.UpdateValue(toState)

	case characteristic.CurrentDoorStateOpen,
		characteristic.CurrentDoorStateClosed:

		log.Println("updateState - [ 0, 1 ]:", toState)
		// Make sure we only update target status if CurrentDoorState is Stopped.
		if oldState != characteristic.CurrentDoorStateStopped {
			log.Println("upateStatus updates targetdoorstate")
			acc.TargetDoorState.UpdateValue(toState)
		}
		// Also, make sure to update TargetDoorState before CurrentDoorState. There
		// are some funky situations where updating TargetDoorState after
		// CurrentDoorState causes a new transition to happen.
		acc.CurrentDoorState.UpdateValue(toState)
	}
}

func startHomeKitServer() {
	hcConfig := hc.Config{
		Pin:         config.Accessory.Pin,
		Port:        config.Accessory.Port,
		SetupId:     config.Accessory.SetupId,
		StoragePath: config.Accessory.StoragePath,
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
