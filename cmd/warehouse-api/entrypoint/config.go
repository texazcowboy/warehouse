package entrypoint

import (
	"github.com/texazcowboy/warehouse/internal/foundation/database"
	"github.com/texazcowboy/warehouse/internal/foundation/parser"
	"github.com/texazcowboy/warehouse/internal/foundation/web"
)

type ApplicationConfig struct {
	ServerCfg   web.WConfig       `yaml:"server"`
	DatabaseCfg database.DBConfig `yaml:"database"`
}

func (c *ApplicationConfig) Read(path *string) {
	// validation tbd
	parser.ParseFile(c, *path)
}
