package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type storage interface {
	Save(settings Settings) error
	Read() *Settings
}

const (
	folder   = "pomogoro"
	filename = "settings.json"
)

var saved *Settings = nil

type jsonStorage struct {
	storage
}

func newStorage() storage {
	return &jsonStorage{}
}

func getPath() string {
	path, _ := os.UserConfigDir()

	return filepath.Join(path, folder)
}

func getFullPath() string {
	fullPath := filepath.Join(getPath(), filename)

	return fullPath
}

func (s *jsonStorage) Save(settings Settings) error {
	saved = &settings

	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}

	test := filepath.Dir(getPath())
	fmt.Printf(test)

	err = os.MkdirAll(getPath(), 0700)
	if err != nil {
		return err
	}

	err = os.WriteFile(getFullPath(), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *jsonStorage) Read() *Settings {
	file, err := os.ReadFile(getFullPath())
	if err != nil {
		return saved
	}

	err = json.Unmarshal(file, &saved)
	if err != nil {
		return saved
	}

	return saved
}
