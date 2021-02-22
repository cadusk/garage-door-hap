package main

import (
        "github.com/brutella/hc/accessory"
        "github.com/brutella/hc/characteristic"
        "github.com/brutella/hc/service"
)

const (
        AccessoryName         string = "GarageDoor"
        AccessoryManufacturer string = "OsCardoso, Inc."
        AccessoryModel        string = "GDRPI"
)

type GarageDoor struct {
        *accessory.Accessory
        *service.GarageDoorOpener
}

func NewGarageDoor(serialNumber, firmwareVersion string) *GarageDoor {
        acc := GarageDoor{}

        info := accessory.Info{
                Name:             AccessoryName,
                SerialNumber:     serialNumber,
                Manufacturer:     AccessoryManufacturer,
                Model:            AccessoryModel,
                FirmwareRevision: firmwareVersion,
        }
        acc.Accessory = accessory.New(info, accessory.TypeGarageDoorOpener)
        acc.GarageDoorOpener = service.NewGarageDoorOpener()

        acc.ObstructionDetected.SetValue(false)

        acc.CurrentDoorState.SetMinValue(characteristic.CurrentDoorStateOpen)
        acc.CurrentDoorState.SetMaxValue(characteristic.CurrentDoorStateStopped)
        acc.CurrentDoorState.SetValue(characteristic.CurrentDoorStateClosed)

        acc.TargetDoorState.SetMinValue(characteristic.TargetDoorStateOpen)
        acc.TargetDoorState.SetMaxValue(characteristic.TargetDoorStateClosed)
        acc.TargetDoorState.SetValue(characteristic.TargetDoorStateClosed)

        acc.AddService(acc.GarageDoorOpener.Service)

        return &acc
}
