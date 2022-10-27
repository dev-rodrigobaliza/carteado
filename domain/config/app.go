package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/dev-rodrigobaliza/carteado/errors"
)

type App struct {
	Date      string    `json:"-"`
	MadeBy    string    `json:"-"`
	Name      string    `json:"-"`
	Version   string    `json:"-"`
	Debug     bool      `json:"-"`
	StartedAt time.Time `json:"-"`
	Timezone  string    `json:"timezone,omitempty"`
	Database  *Database `json:"database,omitempty"`
	HTTP      *HTTP     `json:"http,omitempty"`
	Security  *Security `json:"security,omitempty"`
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
		return errors.ErrFailedOpenFileConfig
	}

	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)

	err = jsonParser.Decode(a)
	if err != nil {
		return errors.ErrFailedParseFileConfig
	}

	return nil
}
