package main

import (
	"fmt"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/versions"
	"github.com/magiconair/properties"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ensureAndroidSdkSetup() error {
	androidSdkRoot := os.Getenv("ANDROID_SDK_ROOT")
	if androidSdkRoot == "" {
		androidSdkRoot = os.Getenv("ANDROID_HOME")
	}

	if androidSdkRoot == "" {
		log.Infof("Neither ANDROID_SDK_ROOT nor ANDROID_HOME environment variable is defined, skipping Android SDK setup.")
		return nil
	}

	sdkProperties, err := properties.LoadFile(androidSdkRoot+"/tools/source.properties", properties.UTF8)
	if err != nil {
		return err
	}

	currentSdkToolsVersion := sdkProperties.GetString("Pkg.Revision", "0.0.0")
	isCurrentSdkToolsUpToDate, err := versions.IsVersionGreaterOrEqual(currentSdkToolsVersion, "26.0.0")
	if err != nil {
		return err
	}

	sdkManagerPath, err := findSdkManagerPath(androidSdkRoot)
	if err != nil {
		return err
	}

	if !isCurrentSdkToolsUpToDate {
		log.Infof("Current Android SDK version: %s is lower than 26. Updating...", currentSdkToolsVersion)
		updateCommand := fmt.Sprintf("yes|%s --update", sdkManagerPath)
		if err := command.RunBashCommand(updateCommand); err != nil {
			return err
		}
	}

	licenseAcceptCommand := fmt.Sprintf("yes|%s --licenses", sdkManagerPath)
	err = command.RunBashCommand(licenseAcceptCommand)
	if err != nil {
		return err
	}

	cmdlineToolsInstallCommand := fmt.Sprintf("yes|%s \"cmdline-tools;7.0\"", sdkManagerPath)
	return command.RunBashCommand(cmdlineToolsInstallCommand)
}

func findSdkManagerPath(androidSdkRoot string) (string, error) {
	var sdkmanagerPath string
	var err = filepath.WalkDir(androidSdkRoot, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, "/sdkmanager") && strings.Contains(path, "/cmdline-tools/") && d.Type().IsRegular() {
			sdkmanagerPath = path
		}
		return err
	})
	return sdkmanagerPath, err
}
