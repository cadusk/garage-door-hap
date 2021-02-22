package main

import (
        "log"
        "time"

        "periph.io/x/conn/v3/gpio"
        "periph.io/x/conn/v3/i2c"
        "periph.io/x/conn/v3/i2c/i2creg"
        "periph.io/x/devices/v3/ep0099"
        "periph.io/x/host/v3"
        "periph.io/x/host/v3/rpi"
)

const clickDuration = 1500 * time.Millisecond

type Device struct {
        bus i2c.BusCloser
        *ep0099.Dev
        sensorOpen, sensorClosed gpio.PinIn
}

func NewDevice(addr uint16) *Device {
        // Initializes host to manage bus and devices
        if _, err := host.Init(); err != nil {
                log.Fatal(err)
        }

        // Opens default bus
        bus, err := i2creg.Open("")
        if err != nil {
                log.Fatal(err)
        }

        sensorOpen := rpi.P1_16 // GPIO23
        if err := sensorOpen.In(gpio.PullUp, gpio.BothEdges); err != nil {
                log.Fatal(err)
        }

        sensorClosed := rpi.P1_18 // GPIO24
        if err := sensorClosed.In(gpio.PullUp, gpio.BothEdges); err != nil {
                log.Fatal(err)
        }

        // Initializes device using current I2C bus and device address.
        // Address should be provided as configured on the board's DIP switches.
        ep0099, err := ep0099.New(bus, addr)
        if err != nil {
                log.Fatal(err)
        }

        return &Device{
                Dev:          ep0099,
                bus:          bus,
                sensorOpen:   sensorOpen,
                sensorClosed: sensorClosed,
        }
}

func (c *Device) Halt() {
        c.bus.Close()
        c.Dev.Halt()
}

func (c *Device) Click(channel uint8) error {
        if err := c.On(channel); err != nil {
                return err
        }

        time.Sleep(clickDuration)

        if err := c.Off(channel); err != nil {
                return err
        }
        return nil
}

func (c *Device) sensorOpenStatus() bool {
        return c.sensorOpen.Read() == gpio.Low
}

func (c *Device) sensorClosedStatus() bool {
        return c.sensorClosed.Read() == gpio.Low
}
