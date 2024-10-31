package settings

type storage interface {
	Save(settings Settings) error
	Read() *Settings
}

var saved *Settings = nil

type jsonStorage struct {
	storage
}

func newStorage() storage {
	return &jsonStorage{}
}

func (s *jsonStorage) Save(settings Settings) error {
	saved = &settings

	return nil
}

func (s *jsonStorage) Read() *Settings {
	return saved
}
