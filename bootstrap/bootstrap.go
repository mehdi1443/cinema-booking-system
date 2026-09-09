package bootstrap

import (
	"fmt"
)

type Config struct {
	Constant    *Constant    `json:"constant"`
	Environment *Environment `json:"environment"`
}

func Run() *Config {
	environment := NewEnvironment()
	constant := NewConstant()
	constant.User = SuperUser{
		Phone:     environment.SuperAdminPhone,
		Email:     environment.SuperAdminEmail,
		Password:  environment.SuperAdminPassword,
		FirstName: environment.SuperAdminFirstName,
		LastName:  environment.SuperAdminLastName,
		Wallet:    environment.SuperAdminWallet,
	}
	return &Config{
		Constant:    constant,
		Environment: environment,
	}
}

func (c *Config) String() string {
	if c == nil {
		return "<nil>"
	}

	return fmt.Sprintf("Config{\n  Constant: %s,\n  Environment: %s\n}", c.Constant, c.Environment)
}
