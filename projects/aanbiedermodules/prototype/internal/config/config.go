package config


import (
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	YAML_FILE  = "config.yaml"
	ENV_PREFIX = ""
	ENV_DELIM  = "_"
	CFG_DELIM  = "."
)

func InitializeConfig() Config {
	k := koanf.New(CFG_DELIM)

	// Load config.yaml file configuration
	_ = k.Load(file.Provider(YAML_FILE), yaml.Parser())

	// Load environment configuration
	_ = k.Load(envProvider(), nil)

	cfg := defaultConfig()

	// Unmarshall koanf environment variables into the cfg struct
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("Error while unmarshalling config: %v", err)
	}

	return cfg
}

func envProvider() *env.Env {
	return env.ProviderWithValue(ENV_PREFIX, CFG_DELIM, func(s string, v string) (string, interface{}) {
		key := strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, ENV_PREFIX)), ENV_DELIM, CFG_DELIM)

		// Otherwise, return the plain string.
		return key, v
	})
}
