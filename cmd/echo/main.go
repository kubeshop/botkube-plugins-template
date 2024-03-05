package main

import (
	"context"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/hashicorp/go-plugin"
	"github.com/kubeshop/botkube/pkg/api"
	"github.com/kubeshop/botkube/pkg/api/executor"
	pluginx "github.com/kubeshop/botkube/pkg/plugin"
)

const description = "Echo sends back the command that was specified."

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
		JSONSchema: api.JSONSchema{
			Value: heredoc.Docf(`{
			  "$schema": "http://json-schema.org/draft-04/schema#",
			  "title": "echo",
			  "description": "%s",
			  "type": "object",
			  "properties": {
			    "formatOptions": {
			      "description": "Options to format echoed string",
			      "type": "array",
			      "items": {
			        "type": "string",
			        "enum": [ "bold", "italic" ]
			      }
			    }
			  },
			  "additionalProperties": false
			}`, description),
		},
	}, nil
}

// Execute returns a given command as a response.
//
//nolint:gocritic  //hugeParam: in is heavy (80 bytes); consider passing it by pointer
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
		Message: api.NewCodeBlockMessage(response, true),
	}, nil
}

func (EchoExecutor) Help(context.Context) (api.Message, error) {
	btnBuilder := api.NewMessageButtonBuilder()
	return api.Message{
		Sections: []api.Section{
			{
				Base: api.Base{
					Header:      "Run `echo` commands",
					Description: description,
				},
				Buttons: []api.Button{
					btnBuilder.ForCommandWithDescCmd("Run", "echo 'hello world'"),
				},
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
