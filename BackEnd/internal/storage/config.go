package storage

type Config struct {
	// Строка подключения к бд
	DatabaseURI string
}

func NewConfig(db string) *Config {
	return &Config{
		DatabaseURI: db,
	}
}
