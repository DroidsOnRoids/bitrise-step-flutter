package main

import (
	"testing"
	"github.com/stretchr/testify/require"
	"os"
)

func TestValidateConfigsNoVersion(t *testing.T) {
	configs := ConfigsModel{
		WorkingDir: "",
		Commands:   []string{"test"},
	}
	require.Error(t, configs.validate())
}

func TestValidateConfigsNoWorkingDir(t *testing.T) {
	configs := ConfigsModel{
		Version:  "1",
		Commands: []string{"test"},
	}
	require.Error(t, configs.validate())
}

func TestValidateConfigsNoCommands(t *testing.T) {
	configs := ConfigsModel{
		Version:    "1",
		WorkingDir: ".",
		Commands:   []string{""},
	}
	require.Error(t, configs.validate())
}

func TestCreateConfigsModelFromEnvsVersion(t *testing.T) {
	err := os.Setenv("version", "123")
	require.NoError(t, err)

	configs := createConfigsModelFromEnvs()

	require.Equal(t, "123", configs.Version)
}

func TestCreateConfigsModelFromEnvsWorkingDir(t *testing.T) {
	err := os.Setenv("working_dir", "/tmp")
	require.NoError(t, err)

	configs := createConfigsModelFromEnvs()

	require.Equal(t, "/tmp", configs.WorkingDir)
}

func TestCreateConfigsModelFromEnvsCommands(t *testing.T) {
	err := os.Setenv("commands", "doctor|test")
	require.NoError(t, err)

	configs := createConfigsModelFromEnvs()

	require.Len(t, configs.Commands, 2)
	require.Equal(t, "doctor", configs.Commands[0])
	require.Equal(t, "test", configs.Commands[1])
}

func TestCreateConfigsModelFromEnvsEmptyCommands(t *testing.T) {
	err := os.Setenv("commands", "||doctor||test||")
	require.NoError(t, err)

	configs := createConfigsModelFromEnvs()

	require.Len(t, configs.Commands, 2)
	require.Equal(t, "doctor", configs.Commands[0])
	require.Equal(t, "test", configs.Commands[1])
}
