package main

import (
	"flag"
	"io/ioutil"

	dupont "github.com/thomas-maurice/dupont/pkg"
	"github.com/thomas-maurice/dupont/pkg/types"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "Configuration file to use")
}

func loadConfig(filePath string) (*types.Config, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config types.Config
	err = yaml.Unmarshal(b, &config)
	return &config, err
}

func main() {
	flag.Parse()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	cfg, err := loadConfig(configFile)
	if err != nil {
		log.Fatal("could not load config", zap.Error(err))
	}

	log.Info("loaded config", zap.Any("config", cfg))

	err = dupont.ApplyConfiguration(log, cfg)
	if err != nil {
		log.Fatal("could not apply config", zap.Error(err))
	}
}
