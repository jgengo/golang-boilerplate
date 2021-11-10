package config

import (
	"io/ioutil"

	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerPort int
	DSN        string `yaml:"dsn" env:"DSN"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.DSN, validation.Required),
	)

}

func Load(file string, port int) (*Config, error) {
	c := Config{
		ServerPort: port,
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	if err = c.Validate(); err != nil {
		return nil, err
	}

	return &c, err

}
