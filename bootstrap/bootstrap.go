package bootstrap

type Config struct {
	Constant    *Constant
	Environment *Environment
}

func Run() *Config {
	return &Config{
		Constant:    NewConstant(),
		Environment: NewEnvironment(),
	}
}
