package main

import (
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCreateConfigsModelFromEnvs(t *testing.T) {
	err := os.Setenv("commands", "||doctor||test||||")
	require.NoError(t, err)

	err = os.Setenv("working_dir", os.TempDir())
	require.NoError(t, err)

	err = os.Setenv("version", "1")
	require.NoError(t, err)

	var config Config
	require.NoError(t, stepconf.Parse(&config))

	config.stripEmptyCommands()
	stepconf.Print(config)

	require.Len(t, config.Commands, 2)
	require.Equal(t, "doctor", config.Commands[0])
	require.Equal(t, "test", config.Commands[1])
}
