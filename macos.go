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

	if err := command.RunCommand("python", "-c", "import six"); err != nil {
		log.Infof("six Python module not found, installing it with easy_install...")
		return command.RunCommand("sudo", "easy_install-3.7", "six")
	}

	return nil
}
