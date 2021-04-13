package appsettings

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Utilitário para obtenção de configurações da aplicação

// AppSettings representa a estrutura que o arquivo '.json' de configurações deve possuir
type AppSettings struct {
	ServerInfo         ServerInfo         `json:"server_info"`
	Environment        Environment        `json:"environment"`
	StorageCredentials StorageCredentials `json:"storage_credentials"`
}

// ServerInfo configurações do servidor HTTP
type ServerInfo struct {
	Address string `json:"address"`
}

type StorageCredentials struct {
	DBHost     string `json:"db_host"`
	DBName     string `json:"db_name"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBPort     int    `json:"db_port"`
}

type Environment struct {
	AppMode      string `json:"app_mode"`
	IsProduction bool   `json:"is_production"`
}

// readFromFile lê as configurações de um arquivo '.json'
func readFromFile(settingsFilePath string) (*AppSettings, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("Failed to get os.Executable : %w", err)
	}

	exPath := filepath.Dir(ex)
	settingsFileAbs, err := filepath.Abs(exPath + "/" + settingsFilePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to get settings File Abs : %w", err)
	}

	var settings AppSettings

	file, err := os.Open(settingsFileAbs)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s : %w", settingsFileAbs, err)
	}

	defer file.Close()

	if err := json.NewDecoder(bufio.NewReader(file)).Decode(&settings); err != nil {
		return nil, fmt.Errorf("Failed to decode %s : %w", settingsFileAbs, err)
	}

	return &settings, nil
}

func readEnv() *Environment {
	return &Environment{
		AppMode:      os.Getenv("APP_MODE"),
		IsProduction: os.Getenv("APP_MODE") == "production",
	}
}

func NewAppSettings(settingsFilePath string) (*AppSettings, error) {
	appSettings, err := readFromFile(settingsFilePath)
	if err != nil {
		return nil, err
	}

	appSettings.Environment = *readEnv()

	return appSettings, nil
}
