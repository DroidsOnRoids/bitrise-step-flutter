package main

// Config ...
type Config struct {
	Version    string   `env:"version,required"`
	WorkingDir string   `env:"working_dir,dir"`
	Commands   []string `env:"commands"`
}

func (config *Config) stripEmptyCommands() {
	var strippedCommands []string
	for _, command := range config.Commands {
		if command != "" {
			strippedCommands = append(strippedCommands, command)
		}
	}
	config.Commands = strippedCommands
}