package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type ProfileConfig struct {
	AnsiColorsEnabled string `json:"ansi-colors-enabled"`
	OAuthClientID     string `json:"oauth-client-id"`
	Output            string `json:"output"`
	SASEndpoint       string `json:"sas-endpoint"`
}

type Config map[string]ProfileConfig

func LoadConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, ".sas", "config.json")
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not parse config file: %v", err)
	}

	return config, nil
}

func main() {

	profile := flag.String("profile", "default", "Specify the profile name")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	profileConfig, exists := config[*profile]
	if !exists {
		fmt.Printf("Profile '%s' not found in config\n", *profile)
		os.Exit(1)
	}

	command := flag.Arg(0)

	log := func(msg string) {
		if *verbose {
			fmt.Println("[DEBUG]", msg)
		}
	}

	log(fmt.Sprintf("Loaded profile: %s", *profile))
	log(fmt.Sprintf("SAS Endpoint: %s", profileConfig.SASEndpoint))

	switch command {
	case "auth":
		log(fmt.Sprintf("Using profile: %s", *profile))
		fmt.Printf("Running auth with profile %s\n", *profile)
		err := handleAuth(*profile, profileConfig)
		if err != nil {
			fmt.Println("Authentication Error:", err)
		}
	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
