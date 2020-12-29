package database

type DBConfig struct {
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
}
