package main

import "time"

type Config struct {
        Accessory struct {
                SerialNumber string        `fig:"serial_number" validate:"required"`
                Pin          string        `fig:"pin" default:"00102003"`
                Timeout      time.Duration `fig:"timeout" default:"12s"`
        } `fig:"accessory"`

        Device struct {
                Address uint16 `fig:"address" default:"0x10"`
                Channel uint8  `fig:"channel" default:"1"`
        } `fig:"device"`

        Server struct {
                StoragePath string `fig:"storage_path" default:"./data"`
                Port        string `fig:"port" default:"57658"`
        } `fig:"server"`
}
