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

func LoadTopology(log *zap.Logger, filePath string, format string) (*types.Topology, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	switch format {
	case "yaml":
		return getYamlTopologyFromBytes(b)
	case "hcl":
		return getHCLTopologyFromBytes(log, b, filePath)
	default:
		return nil, fmt.Errorf("unknown configuration format %s", format)
	}
}

func getYamlTopologyFromBytes(b []byte) (*types.Topology, error) {
	var config types.Topology
	err := yaml.Unmarshal(b, &config)
	return &config, err
}

func getHCLTopologyFromBytes(log *zap.Logger, b []byte, fileName string) (*types.Topology, error) {
	var cfg types.Topology
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
