package main

// Config ...
type Config struct {
	Version    string   `env:"version,required"`
	WorkingDir string   `env:"working_dir,dir"`
	Commands   []string `env:"commands"`
}

func (config *Config) stripEmptyCommands() {
	var commands []string
	for _, pth := range config.Commands {
		if pth != "" {
			commands = append(commands, pth)
		}
	}
	config.Commands = commands
}