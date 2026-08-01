package bootstrap

import (
	"fmt"
)

type Config struct {
	Constant    *Constant    `json:"constant"`
	Environment *Environment `json:"environment"`
}

func Run() *Config {
	return &Config{
		Constant:    NewConstant(),
		Environment: NewEnvironment(),
	}
}

// TODO : Declare String() Func for Debug
func (c *Config) String() string {
	if c == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Config{\n  Constant: %s,\n  Environment: %s\n}", c.Constant, c.Environment)
}
