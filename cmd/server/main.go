package main

import (
	"flag"
	"fmt"

	"github.com/jgengo/golang-boilerplate/internal/config"
)

var Version = "0.0.1"

var flagConfig = flag.String("config", "./config/dev.yml", "path to the config file")

func main() {
	flag.Parse()

	cfg, err := config.Load(*flagConfig)

	fmt.Println(cfg)
	fmt.Println(err)

}
