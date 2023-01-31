package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/kubeshop/botkube/pkg/api"
	"github.com/kubeshop/botkube/pkg/api/executor"
	"github.com/kubeshop/botkube/pkg/bot/interactive"
	"github.com/kubeshop/botkube/pkg/pluginx"
)

const (
	description = "Echo sends back the command that was specified."
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
		Description: description,
	}, nil
}

// Execute returns a given command as a response.
func (EchoExecutor) Execute(_ context.Context, in executor.ExecuteInput) (executor.ExecuteOutput, error) {
	var cfg Config
	err := pluginx.MergeExecutorConfigs(in.Configs, &cfg)
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

// Help returns help message
func (EchoExecutor) Help(context.Context) (interactive.Message, error) {
	return interactive.Message{
		Base: interactive.Base{
			Body: interactive.Body{
				Plaintext: description,
			},
		},
	}, nil
}

func main() {
	executor.Serve(map[string]plugin.Plugin{
		"echo": &executor.Plugin{
			Executor: &EchoExecutor{},
		},
	})
}
