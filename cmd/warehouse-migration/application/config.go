package application

import (
	"github.com/texazcowboy/warehouse/internal/foundation/database"
	"github.com/texazcowboy/warehouse/internal/foundation/logger"
	"github.com/texazcowboy/warehouse/internal/foundation/parser"
)

type Config struct {
	*database.DBConfig `yaml:"database" validate:"required"`
	*logger.LConfig    `yaml:"log" validate:"required"`
}

func (c *Config) Read(path string) error {
	return parser.ParseFile(c, path)
}
