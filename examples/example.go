package main

import (
	"fmt"
	"github.com/koofr/envigo"
)

type Config struct {
	Http    ConfigHttp
	Logging ConfigLogging
	Debug   bool
}

type ConfigHttp struct {
	Port int
}

type ConfigLogging struct {
	Level string
}

func main() {
	config := &Config{}

	_ = envigo.Envigo(config, "", envigo.EnvironGetter())

	fmt.Println(config)
}
