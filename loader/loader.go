package loader

import (
	"encoding/json"
	"io/ioutil"
)

// Config interface
type Config interface {
	Load(path string, v interface{}) error
}

// LoadConfig method
func LoadConfig(c Config, p string, v interface{}) error {
	return c.Load(p, v)
}

type config struct{}

func (config) Load(p string, v interface{}) error {
	data, err := ioutil.ReadFile(p)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)

	if err != nil {
		return err
	}

	return nil
}

// NewConfig method
func NewConfig() Config {
	return config{}
}
