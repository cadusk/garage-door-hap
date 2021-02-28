package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kkyr/fig"
)

const defaultConfigFilename string = "config.yaml"

type Config struct {
	Accessory struct {
		SerialNumber string        `fig:"serial_number" validate:"required"`
		Pin          string        `fig:"pin" default:"00102003"`
		Timeout      time.Duration `fig:"timeout" default:"12s"`
		StoragePath  string        `fig:"storage_path"`
		SetupId      string        `fig:"setup_id"`
		Port         string        `fig:"port" default:"57658"`
	} `fig:"accessory"`

	Device struct {
		Address uint16 `fig:"address" default:"0x10"`
		Channel uint8  `fig:"channel" default:"1"`
	} `fig:"device"`
}

func loadConfig(config *Config) error {

	dir, err := findExecutableDirectory()
	if err != nil {
		return err
	}

	options := []fig.Option{
		fig.File(defaultConfigFilename),
		fig.Dirs(dir),
	}

	err = fig.Load(config, options[:]...)
	if err == nil {
		log.Printf("%+v\n", config)
	}

	// validations
	ensureStoragePath(config)

	return err
}

func findExecutableDirectory() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func ensureStoragePath(config *Config) error {
	if config.Accessory.StoragePath != "" {
		return nil
	}

	executablePath, err := findExecutableDirectory()
	if err != nil {
		return err
	}

	config.Accessory.StoragePath = filepath.Join(executablePath, "data")
	return nil
}
