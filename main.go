package main

import (
	"fmt"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/git"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	var config Config
	if err := stepconf.Parse(&config); err != nil {
		log.Errorf("Configuration error: %s\n", err)
		os.Exit(7)
	}
	stepconf.Print(config)

	if err := ensureAndroidSdkSetup(); err != nil {
		log.Errorf("Could not setup Android SDK, error: %s", err)
		os.Exit(6)
	}

	flutterSdkDir, err := getSdkDestinationDir()
	if err != nil {
		log.Errorf("Could not calculate Flutter SDK destination directory path, error: %s", err)
		os.Exit(5)
	}

	flutterSdkExists, err := pathutil.IsDirExists(flutterSdkDir)
	if err != nil {
		log.Errorf("Could not check if Flutter SDK is installed, error: %s", err)
		os.Exit(1)
	}

	if !flutterSdkExists {
		log.Infof("Extracting Flutter SDK to %s", flutterSdkDir)

		if err := downloadAndExtractReleaseSdk(config.Version, flutterSdkDir); err != nil {
			log.Infof("Version %s not found in releases, trying snapshot.", config.Version)

			if err := downloadAndExtractSnapshotSdk(config.Version, flutterSdkDir); err != nil {
				log.Errorf("Could not extract Flutter SDK, error: %s", err)
				os.Exit(2)
			}
		}
	} else {
		log.Infof("Flutter SDK directory already exists, skipping installation.")

		flutterVersion, err := fileutil.ReadStringFromFile(flutterSdkDir + "/version")
		if err != nil {
			log.Warnf("Could not determine installed Flutter version, error: %s", err)
		} else if flutterVersion != config.Version {
			log.Warnf("Already installed Flutter version %s will be used instead of requested version %s ", flutterVersion, config.Version)
		}
	}

	for _, flutterCommand := range config.Commands {
		log.Infof("Executing Flutter command: %s", flutterCommand)

		flutterExecutablePath := filepath.Join(flutterSdkDir, "bin/flutter")
		bashCommand := fmt.Sprintf("%s %s", flutterExecutablePath, flutterCommand)
		err := command.RunCommandInDir(config.WorkingDir, "bash", "-c", bashCommand)
		if err != nil {
			log.Errorf("Flutter invocation failed, error: %s", err)
			os.Exit(3)
		}
	}
}

func downloadAndExtractReleaseSdk(flutterVersion, flutterSdkDestinationDir string) error {
	flutterSdkSourceURL := getFlutterSdkWithChannelSourceURL(flutterVersion)

	flutterSdkParentDir := filepath.Join(flutterSdkDestinationDir, "..")

	if runtime.GOOS == "darwin" {
		return command.DownloadAndUnZIP(flutterSdkSourceURL, flutterSdkParentDir)
	} else if runtime.GOOS == "linux" {
		return downloadAndUnTarXZ(flutterSdkSourceURL, flutterSdkParentDir)
	} else {
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func getFlutterSdkWithChannelSourceURL(flutterVersion string) string {
	channels := [2]string{"stable", "beta"}
	var url string
	for _, channel := range channels {
		url = getFlutterSdkSourceURL(flutterVersion, channel)
		if sdkFileExists(url) {
			break
		}
	}
	return url
}

func sdkFileExists(url string) bool {
	response, err := http.Head(url)
	return err == nil && response.StatusCode == 200
}

func getFlutterSdkSourceURL(flutterVersion string, channel string) string {
	return fmt.Sprintf(
		"https://storage.googleapis.com/flutter_infra_release/releases/%s/%s/flutter_%s_%s%s-%s.%s",
		channel,
		getFlutterPlatform(),
		getFlutterPlatform(),
		getArchitecture(),
		flutterVersion,
		channel,
		getArchiveExtension())
}

func getArchitecture() string {
	if runtime.GOARCH == "arm64" {
		return "arm64_"
	}
	return ""
}

func downloadAndExtractSnapshotSdk(flutterVersion, flutterSdkDestinationDir string) error {

	flutterVersion = normalizeFlutterVersion(flutterVersion)
	gitRepo, err := git.New(flutterSdkDestinationDir)
	if err != nil {
		return err
	}

	return gitRepo.CloneTagOrBranch("https://github.com/flutter/flutter.git", flutterVersion).Run()
}
