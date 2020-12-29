package web

type WConfig struct {
	Port string `yaml:"port" validate:"required,numeric"`
	// ...
}
