package main

import (
	"os"
	"fmt"
	"strings"
	"errors"
	"github.com/bitrise-io/go-utils/log"
)

// ConfigsModel ...
type ConfigsModel struct {
	Version    string
	WorkingDir string
	Commands   []string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		Version:    os.Getenv("version"),
		WorkingDir: os.Getenv("working_dir"),
		Commands:   strings.Split(os.Getenv("commands"), "|"),
	}
}

func (configs ConfigsModel) dump() {
	fmt.Println()
	log.Infof("Configs:")
	log.Printf(" - Flutter version (hidden): %s", configs.Version)
	log.Printf(" - WorkingDir: %s", configs.WorkingDir)
	log.Printf(" - Commands: %s", configs.Commands)
}

func (configs ConfigsModel) validate() error {
	if configs.Version == "" {
		return errors.New("empty Flutter version specified")
	}

	if configs.WorkingDir == "" {
		return errors.New("empty WorkingDir specified")
	}

	for _, command := range configs.Commands {
		if command == "" {
			return errors.New("empty Flutter command specified")
		}
	}

	return nil
}
