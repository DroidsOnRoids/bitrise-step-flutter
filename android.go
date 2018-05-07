package main

import (
	"os"
	"github.com/bitrise-io/go-utils/log"
	"github.com/magiconair/properties"
	"github.com/bitrise-io/go-utils/versions"
	"path"
	"fmt"
	"github.com/bitrise-io/go-utils/command"
)

func ensureAndroidSdkSetup() error {
	androidHome := os.Getenv("ANDROID_HOME")

	if androidHome == "" {
		log.Infof("ANDROID_HOME environment variable not defined, skipping Android SDK setup.")
		return nil
	}

	sdkProperties, err := properties.LoadFile("${ANDROID_HOME}/tools/source.properties", properties.UTF8)
	if err != nil {
		return err
	}

	currentSdkToolsVersion := sdkProperties.GetString("Pkg.Revision", "0.0.0")
	isCurrentSdkToolsUpToDate, err := versions.IsVersionGreaterOrEqual(currentSdkToolsVersion, "26.0.0")
	if err != nil {
		return err
	}

	sdkManagerPath := path.Join(androidHome, "tools/bin/sdkmanager")

	if !isCurrentSdkToolsUpToDate {
		log.Infof("Current Android SDK version: %s is lower than 26. Updating...", currentSdkToolsVersion)
		updateCommand := fmt.Sprintf("yes|%s --update", sdkManagerPath)
		if err := command.RunBashCommand(updateCommand); err != nil {
			return err
		}
	}

	licenseAcceptCommand := fmt.Sprintf("yes|%s --licenses", sdkManagerPath)
	return command.RunBashCommand(licenseAcceptCommand)
}