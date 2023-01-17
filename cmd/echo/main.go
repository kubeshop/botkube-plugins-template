package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/kubeshop/botkube/pkg/api"
	"github.com/kubeshop/botkube/pkg/api/executor"
	"gopkg.in/yaml.v3"
)

// version is set via ldflags by GoReleaser.
var version = "dev"

// Config holds the executor configuration.
type Config struct {
	TransformResponseToUpperCase *bool `yaml:"transformResponseToUpperCase,omitempty"`
}

// EchoExecutor implements the Botkube executor plugin interface.
type EchoExecutor struct{}

// Metadata returns details about the Echo plugin.
func (EchoExecutor) Metadata(context.Context) (api.MetadataOutput, error) {
	return api.MetadataOutput{
		Version:     version,
		Description: "Echo sends back the command that was specified.",
	}, nil
}

// Execute returns a given command as a response.
func (EchoExecutor) Execute(_ context.Context, in executor.ExecuteInput) (executor.ExecuteOutput, error) {
	cfg, err := mergeConfigs(in.Configs)
	if err != nil {
		return executor.ExecuteOutput{}, err
	}

	response := in.Command
	if cfg.TransformResponseToUpperCase != nil && *cfg.TransformResponseToUpperCase {
		response = strings.ToUpper(response)
	}

	return executor.ExecuteOutput{
		Data: fmt.Sprintf("Echo: %s", response),
	}, nil
}

func main() {
	executor.Serve(map[string]plugin.Plugin{
		"echo": &executor.Plugin{
			Executor: &EchoExecutor{},
		},
	})
}

// mergeConfigs merges all input configuration. In our case we don't have complex merge strategy,
// the last one that was specified wins :)
func mergeConfigs(configs []*executor.Config) (Config, error) {
	finalCfg := Config{}
	for _, inputCfg := range configs {
		var cfg Config
		err := yaml.Unmarshal(inputCfg.RawYAML, &cfg)
		if err != nil {
			return Config{}, fmt.Errorf("while unmarshalling YAML config: %w", err)
		}
		if cfg.TransformResponseToUpperCase == nil {
			continue
		}
		finalCfg.TransformResponseToUpperCase = cfg.TransformResponseToUpperCase
	}

	return finalCfg, nil
}
