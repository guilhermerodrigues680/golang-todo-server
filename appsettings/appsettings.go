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
	ServerInfo ServerInfo `json:"server_info"`
}

// ServerInfo configurações do servidor HTTP
type ServerInfo struct {
	Address string `json:"address"`
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

func NewAppSettings(settingsFilePath string) (*AppSettings, error) {
	return readFromFile(settingsFilePath)
}
