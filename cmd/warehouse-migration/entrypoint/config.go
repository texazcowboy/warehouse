package entrypoint

import (
	"github.com/texazcowboy/warehouse/internal/foundation/database"
	"github.com/texazcowboy/warehouse/internal/foundation/logger"
	"github.com/texazcowboy/warehouse/internal/foundation/parser"
)

type ApplicationConfig struct {
	*database.DBConfig `yaml:"database"`
	*logger.LConfig    `yaml:"logger"`
}

func (c *ApplicationConfig) Read(path *string) error {
	return parser.ParseFile(c, *path)
}