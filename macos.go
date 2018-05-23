package main

import (
	"github.com/bitrise-io/go-utils/command"
	"runtime"
	"github.com/bitrise-io/go-utils/log"
)

func ensureMacOSSetup() error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	if err := command.RunCommand("python", "-c", "import six"); err != nil {
		log.Infof("six Python module not found, installing it with easy_install...")
		return command.RunCommand("sudo", "easy_install six")
	}

	return nil
}
