package settings

import (
	"encoding/json"
	"os"
)

type storage interface {
	Save(settings Settings) error
	Read() *Settings
}

const filePath = "settings.json"

var saved *Settings = nil

type jsonStorage struct {
	storage
}

func newStorage() storage {
	return &jsonStorage{}
}

// todo: os.UserConfigDir
func (s *jsonStorage) Save(settings Settings) error {
	saved = &settings

	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *jsonStorage) Read() *Settings {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return saved
	}

	err = json.Unmarshal(file, &saved)
	if err != nil {
		panic(err)
	}

	return saved
}
