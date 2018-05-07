package main

// Config ...
type Config struct {
	Version    string   `env:"version,required"`
	WorkingDir string   `env:"working_dir,dir"`
	Commands   []string `env:"commands,required"`
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