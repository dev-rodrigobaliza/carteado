package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type App struct {
	Date     string    `json:"-"`
	MadeBy   string    `json:"-"`
	Name     string    `json:"-"`
	Version  string    `json:"-"`
	Debug    bool      `json:"-"`
	Timezone string    `json:"timezone,omitempty"`
	Database *Database `json:"database,omitempty"`
	HTTP     *HTTP     `json:"http,omitempty"`
	Security *Security `json:"security,omitempty"`
}

func NewApp(name, version, date, madeBy string, debug bool) *App {
	return &App{
		Date:    date,
		MadeBy:  madeBy,
		Name:    name,
		Version: version,
		Debug:   debug,
	}
}

func (a *App) LoadFromFile(filename string) error {
	configFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open config file (%s)", filename)
	}

	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)

	err = jsonParser.Decode(a)
	if err != nil {
		return fmt.Errorf("failed to parse config file (%s)", filename)
	}

	return nil
}
