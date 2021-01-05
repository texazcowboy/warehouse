package application

import (
	"github.com/texazcowboy/warehouse/internal/foundation/database"
	"github.com/texazcowboy/warehouse/internal/foundation/logger"
	"github.com/texazcowboy/warehouse/internal/foundation/parser"
	"github.com/texazcowboy/warehouse/internal/foundation/web"
)

type ApplicationConfig struct {
	ServerCfg   *web.WConfig       `yaml:"server" validate:"required"`
	DatabaseCfg *database.DBConfig `yaml:"database" validate:"required"`
	LoggerCfg   *logger.LConfig    `yaml:"log" validate:"required"`
}

func (c *ApplicationConfig) Read(path *string) error {
	return parser.ParseFile(c, *path)
}
