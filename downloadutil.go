package main

import (
	"fmt"
	"github.com/bitrise-io/go-utils/log"
	"github.com/blang/semver"
	"github.com/mholt/archiver"
	"io"
	"net/http"
	"os"
)

func downloadAndUnTarXZ(url, dirPath string) error {
	file, err := os.CreateTemp("", "flutter.*.tar.xz")
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("Failed to close temporary file %s, error: %s", file.Name(), err)
		}
		if err := os.Remove(file.Name()); err != nil {
			log.Errorf("Failed to remove temporary file %s, error: %s", file.Name(), err)
		}
	}()

	if err := downloadFile(url, file); err != nil {
		return err
	}

	return archiver.Unarchive(file.Name(), dirPath)
}

func downloadFile(downloadURL string, outFile *os.File) error {
	response, err := http.Get(downloadURL)
	if err != nil {
		return err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Warnf("Failed to close (%s) body", downloadURL)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file from %s, error: %s", downloadURL, response.Status)
	}

	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		return fmt.Errorf("failed to save file %s, error: %s", outFile.Name(), err)
	}

	return nil
}

func normalizeFlutterVersion(flutterVersion string) string {
	if version, err := semver.Parse(flutterVersion); err == nil {
		if version.LT(semver.MustParse("1.17.0")) {
			return "v" + flutterVersion
		}
	}
	return flutterVersion
}
