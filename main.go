package main

import (
	"flag"

	dupont "github.com/thomas-maurice/dupont/pkg"
	"github.com/thomas-maurice/dupont/pkg/config"
	"go.uber.org/zap"
)

var (
	configFile   string
	configFormat string
)

func init() {
	flag.StringVar(&configFile, "config", "config.hcl", "Configuration file to use")
	flag.StringVar(&configFormat, "config-format", "hcl", "Configuration format (hcl/yaml)")
}

func main() {
	flag.Parse()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	cfg, err := config.LoadConfig(log, configFile, configFormat)
	if err != nil {
		log.Fatal("could not load config", zap.Error(err))
	}

	log.Info("loaded config", zap.Any("config", cfg))

	err = dupont.ApplyConfiguration(log, cfg)
	if err != nil {
		log.Fatal("could not apply config", zap.Error(err))
	}
}
