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
	var commands []string
	for _, pth := range strings.Split(os.Getenv("commands"), "|") {
		if pth != "" {
			commands = append(commands, pth)
		}
	}

	return ConfigsModel{
		Version:    os.Getenv("version"),
		WorkingDir: os.Getenv("working_dir"),
		Commands:   commands,
	}
}

func (configs ConfigsModel) dump() {
	fmt.Println()
	log.Infof("Configs:")
	log.Printf(" - Flutter version: %s", configs.Version)
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

	if len(configs.Commands) == 0 {
		return errors.New("no Flutter command specified")
	}

	return nil
}
