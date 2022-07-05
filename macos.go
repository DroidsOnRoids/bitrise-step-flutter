package main

import (
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"runtime"
)

func ensureMacOSSetup() error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	if err := command.RunCommand("python3", "-c", "import six"); err != nil {
		log.Infof("six Python module not found, installing it")
		return command.RunCommand("sudo", "pip3", "install", "six")
	}

	return nil
}
