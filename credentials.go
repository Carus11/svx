package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CredentialItem struct {
	AccessToken  string    `json:"access-token"`
	Expiry       time.Time `json:"expiry"`
	RefreshToken string    `json:"refresh-token"`
}

type Credentials map[string]CredentialItem

func getCredentialsFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".sas", "credentials.json")
}

func LoadCredentials() (Credentials, error) {
	credentialsPath := getCredentialsFilePath()

	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		return Credentials{}, nil
	}

	data, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read credentials file: %v", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("Failed to parse credentials JSON: %v", err)
	}

	return creds, nil
}

func SaveCredentials(creds Credentials) error {
	credentialsPath := getCredentialsFilePath()

	data, err := json.MarshalIndent(creds, "", " ")
	if err != nil {
		return fmt.Errorf("Failed to marshal credentials JSON: %v", err)
	}

	err = os.WriteFile(credentialsPath, data, 0600)
	if err != nil {
		return fmt.Errorf("Failed write credentials file: %v", err)
	}

	return nil

}

func GetToken(profile string) (*CredentialItem, error) {
	creds, err := LoadCredentials()
	if err != nil {
		return nil, err
	}

	credentialItem, exists := creds[profile]

	if !exists {
		return nil, fmt.Errorf("No token found for profile %s", profile)
	}

	return &credentialItem, nil
}

func SaveToken(profile string, credentialItem *CredentialItem) error {
	creds, err := LoadCredentials()
	if err != nil {
		return err
	}

	creds[profile] = *credentialItem

	return SaveCredentials(creds)
}
