package main

import (
	"fmt"
	"github.com/bitrise-io/go-utils/pathutil"
	"path/filepath"
	"runtime"
)

func getSdkDestinationDir() (string, error) {
	if runtime.GOOS == "darwin" {
		return filepath.Join(pathutil.UserHomeDir(), "Library/flutter"), nil
	} else if runtime.GOOS == "linux" {
		return "/opt/flutter", nil
	}
	return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
}

func getFlutterPlatform() string {
	if runtime.GOOS == "darwin" {
		return "macos"
	}
	return runtime.GOOS
}

func getArchiveExtension() string {
	if runtime.GOOS == "linux" {
		return "tar.xz"
	}
	return "zip"
}
