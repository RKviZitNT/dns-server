package config

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

type Config struct {
	Address string
}

// загрузка конфиг файла
func LoadConfig(dir, name string) (*Config, error) {
	nameAndType := strings.Split(name, ".")
	if len(nameAndType) != 2 {
		return nil, fmt.Errorf("invalid config file name: %s", name)
	}

	filePath := path.Join(dir, name)

	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipFlags:         true,
		SkipEnv:           true,
		SkipFiles:         false,
		AllowUnknownFlags: true,
		Files:             []string{filePath},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})

	if err := loader.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("Config loaded successfully")
	return &cfg, nil
}
