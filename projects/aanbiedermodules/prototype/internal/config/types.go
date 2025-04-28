package config

type Config struct {
	Stelselnode Stelselnode        `koanf:"stelselnode"`
	LogLevel    string        `koanf:"loglevel"`
	Port        int           `koanf:"port"`
}

type Stelselnode struct {
	Cert []byte   `koanf:"cert"`
	Key []byte`koanf:"key"`
}
