package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"

	"go.uber.org/zap"

	"github.com/thomas-maurice/dupont/pkg/types"

	"gopkg.in/yaml.v3"
)

func LoadConfig(log *zap.Logger, filePath string, format string) (*types.Config, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	switch format {
	case "yaml":
		return getYamlConfigFromBytes(b)
	case "hcl":
		return getHCLConfigFromBytes(log, b, filePath)
	default:
		return nil, fmt.Errorf("unknown configuration format %s", format)
	}
}

func getYamlConfigFromBytes(b []byte) (*types.Config, error) {
	var config types.Config
	err := yaml.Unmarshal(b, &config)
	return &config, err
}

func getHCLConfigFromBytes(log *zap.Logger, b []byte, fileName string) (*types.Config, error) {
	var cfg types.Config
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL(b, fileName)
	if diags.HasErrors() {
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout,
			parser.Files(),
			78,
			true,
		)
		err := wr.WriteDiagnostics(diags)
		if err != nil {
			log.Error("could not write config diagnostic", zap.Error(err))
		}
		return nil, fmt.Errorf("invalid configutation")
	}

	diags = gohcl.DecodeBody(f.Body, &hcl.EvalContext{}, &cfg)
	if diags.HasErrors() {
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout,
			parser.Files(),
			78,
			true,
		)
		err := wr.WriteDiagnostics(diags)
		if err != nil {
			log.Error("could not write config diagnostic", zap.Error(err))
		}
		return nil, fmt.Errorf("invalid configuration")
	}

	return &cfg, nil
}
