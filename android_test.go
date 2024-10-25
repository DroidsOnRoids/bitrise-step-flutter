package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFindBuildToolsVersion(t *testing.T) {
	androidHome := os.Getenv("ANDROID_HOME")
	require.NotEmpty(t, androidHome)
	properties, err := findBuildToolsVersion(androidHome)

	require.NoError(t, err)
	require.NotEmpty(t, properties)
}
