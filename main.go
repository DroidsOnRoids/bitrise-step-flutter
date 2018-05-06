package main

import (
	"github.com/bitrise-io/go-utils/log"
	"os"
	"github.com/bitrise-io/go-utils/command"
	"fmt"
	"runtime"
	"github.com/bitrise-io/go-utils/pathutil"
	"strings"
	"path/filepath"
)

func main() {
	if err := ensureAndroidSdkSetup(); err != nil {
		log.Errorf("Could not setup Android SDK, error: %s", err)
		os.Exit(6)
	}

	configs := createConfigsModelFromEnvs()

	if err := configs.validate(); err != nil {
		log.Errorf("Could not validate config, error: %s", err)
		os.Exit(4)
	}
	configs.dump()

	flutterSdkDir, err := getSdkDestinationDir()
	if err != nil {
		log.Errorf("Could not validate config, error: %s", err)
		os.Exit(5)
	}

	flutterSdkExists, err := pathutil.IsDirExists(flutterSdkDir)
	if err != nil {
		log.Errorf("Could not check if Flutter SDK is installed, error: %s", err)
		os.Exit(1)
	}

	if !flutterSdkExists {
		if err := extractSdk(configs.Version, flutterSdkDir); err != nil {
			log.Errorf("Could not extract Flutter SDK, error: %s", err)
			os.Exit(2)
		}
	} else {
		log.Infof("Flutter SDK folder already exists, skipping installation.")
	}

	for _, flutterCommand := range configs.Commands {
		log.Infof("Executing Flutter command: %s", flutterCommand)

		flutterExecutablePath := filepath.Join(flutterSdkDir, "bin/flutter")
		bashCommand := fmt.Sprintf("%s %s", flutterExecutablePath, flutterCommand)
		err := command.RunCommandInDir(configs.WorkingDir, "bash", "-c", bashCommand)
		if err != nil {
			log.Errorf("Flutter invocation failed, error: %s", err)
			os.Exit(3)
		}
	}
}

func extractSdk(flutterVersion, flutterSdkDestinationDir string) error {
	log.Infof("Extracting Flutter SDK to %s", flutterSdkDestinationDir)

	versionComponents := strings.Split(flutterVersion, "-")
	channel := versionComponents[len(versionComponents)-1]

	flutterSdkSourceURL := fmt.Sprintf(
		"https://storage.googleapis.com/flutter_infra/releases/%s/%s/flutter_%s_v%s.%s",
		channel,
		getFlutterPlatform(),
		getFlutterPlatform(),
		flutterVersion,
		getArchiveExtension())

	flutterSdkParentDir := filepath.Join(flutterSdkDestinationDir, "..")

	if runtime.GOOS == "darwin" {
		return command.DownloadAndUnZIP(flutterSdkSourceURL, flutterSdkParentDir)
	} else if runtime.GOOS == "linux" {
		return downloadAndUnTarXZ(flutterSdkSourceURL, flutterSdkParentDir)
	} else {
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}