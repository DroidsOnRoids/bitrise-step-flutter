package main

import (
	"fmt"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/versions"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
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

	sdkManagerPath, err := findSdkManagerPath(androidSdkRoot)
	if err != nil {
		return err
	}

	if !isAndroidBuildToolsUpToDate(androidSdkRoot) {
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

func findBuildToolsVersion(root string) (string, error) {
	var matches []string
	versionRegex := regexp.MustCompile(`build-tools/(\d+\.\d+\.\d+)/source.properties`)

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Warnf("failed to access path %s, skipping, error: %s", path, err)
			return filepath.SkipDir
		}
		if versionRegex.MatchString(path) {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no source.properties file found")
	}

	slices.SortFunc(matches, func(a, b string) int {
		versionA := versionRegex.FindStringSubmatch(a)[1]
		versionB := versionRegex.FindStringSubmatch(b)[1]
		compareVersions, err2 := versions.CompareVersions(versionA, versionB)
		if err2 != nil {
			log.Warnf(
				"Failed to compare versions %s and %s, assuming they are equal. Error: %s",
				versionA,
				versionB,
				err2,
			)
			return 0
		}
		return compareVersions
	})

	return versionRegex.FindStringSubmatch(matches[0])[1], nil
}

func isAndroidBuildToolsUpToDate(androidSdkRoot string) bool {
	currentBuildToolsVersion, err := findBuildToolsVersion(androidSdkRoot)
	if err != nil {
		log.Warnf("Failed to determine current Android Build Tools version, assuming it is outdated. Error: %s", err)
		return false
	}
	isCurrentBuildToolsUpToDate, err := versions.IsVersionGreaterOrEqual(currentBuildToolsVersion, "26.0.0")
	if err != nil {
		log.Warnf(
			"Failed to compare current Android Build Tools version %s with reference, assuming it is outdated. Error: %s",
			currentBuildToolsVersion,
			err,
		)
		return false
	}
	if !isCurrentBuildToolsUpToDate {
		log.Infof("Current Android Build Tools version %s is lower than 26.0.0, updating...", currentBuildToolsVersion)
	}
	return isCurrentBuildToolsUpToDate
}
