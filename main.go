package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	dupont "github.com/thomas-maurice/dupont/pkg"
	"github.com/thomas-maurice/dupont/pkg/config"
	"go.uber.org/zap"
)

var (
	configFile   string
	configFormat string
	what         string
)

func init() {
	flag.StringVar(&configFile, "config", "config.hcl", "Configuration file to use")
	flag.StringVar(&configFormat, "config-format", "hcl", "Configuration format (hcl/yaml)")
	flag.StringVar(&what, "what", "apply", "What to do ? (apply/delete/generate)")
}

func main() {
	flag.Parse()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	switch what {
	case "apply":
		cfg, err := config.LoadConfig(log, configFile, configFormat)
		if err != nil {
			log.Fatal("could not load config", zap.Error(err))
		}

		log.Info("loaded config", zap.Any("config", cfg))

		err = dupont.ApplyConfiguration(log, cfg)
		if err != nil {
			log.Fatal("could not apply config", zap.Error(err))
		}
	case "delete":
		cfg, err := config.LoadConfig(log, configFile, configFormat)
		if err != nil {
			log.Fatal("could not load config", zap.Error(err))
		}

		log.Info("loaded config", zap.Any("config", cfg))

		err = dupont.TearDownConfiguration(log, cfg)
		if err != nil {
			log.Fatal("could not tear down config", zap.Error(err))
		}
	case "generate":
		cfg, err := config.LoadTopology(log, configFile, configFormat)
		if err != nil {
			log.Fatal("could not load topology", zap.Error(err))
		}

		log.Info("loaded topololgy", zap.Any("topololgy", cfg))

		topology, err := dupont.GenerateTopology(cfg)
		if err != nil {
			log.Fatal("could not generate topology", zap.Error(err))
		}

		for hostName, hostConfig := range topology {
			hclFile := hclwrite.NewEmptyFile()

			gohcl.EncodeIntoBody(hostConfig, hclFile.Body())

			wr := bytes.NewBuffer(nil)
			wr.WriteString(fmt.Sprintf(`/*
	Auto generated configuration ! You probably should not touch this.

	Generated on: %v
	Topology ID : %s 
*/
`, time.Now(), cfg.ID()))
			hclFile.WriteTo(wr)

			if _, err := os.Stat(cfg.ID()); os.IsNotExist(err) {
				err := os.Mkdir(cfg.ID(), 0700)
				if err != nil {
					log.Fatal("could not create topology dir", zap.Error(err))
				}
			}

			err := ioutil.WriteFile(path.Join(cfg.ID(), hostName+".hcl"), wr.Bytes(), 0600)
			if err != nil {
				log.Fatal("could not create topology config file", zap.Error(err))
			}
		}
	default:
		log.Fatal("unknown operation", zap.String("operation", what))
	}
}
