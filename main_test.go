package main

import (
	"github.com/stretchr/testify/require"
	"runtime"
	"testing"
)

func Test_getFlutterSdkSourceURL(t *testing.T) {
	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			require.Equal(t, getFlutterSdkSourceURL("3.0.1", "stable"), "https://storage.googleapis.com/flutter_infra_release/releases/stable/macos/flutter_macos_arm64_3.0.1-stable.zip")
			require.Equal(t, getFlutterSdkSourceURL("3.1.0", "beta"), "https://storage.googleapis.com/flutter_infra_release/releases/beta/macos/flutter_macos_arm64_3.1.0-beta.zip")
		} else {
			require.Equal(t, getFlutterSdkSourceURL("3.0.1", "stable"), "https://storage.googleapis.com/flutter_infra_release/releases/stable/macos/flutter_macos_3.0.1-stable.zip")
			require.Equal(t, getFlutterSdkSourceURL("3.1.0", "beta"), "https://storage.googleapis.com/flutter_infra_release/releases/beta/macos/flutter_macos_3.1.0-beta.zip")
		}
	} else if runtime.GOOS == "linux" {
		require.Equal(t, getFlutterSdkSourceURL("3.0.1", "stable"), "https://storage.googleapis.com/flutter_infra_release/releases/stable/linux/flutter_linux_3.0.1-stable.tar.xz")
		require.Equal(t, getFlutterSdkSourceURL("3.1.0", "beta"), "https://storage.googleapis.com/flutter_infra_release/releases/beta/linux/flutter_linux_3.1.0-beta.tar.xz")
	}
}
