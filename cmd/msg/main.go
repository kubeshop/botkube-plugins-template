package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/kubeshop/botkube/pkg/api"
	"github.com/kubeshop/botkube/pkg/api/executor"
)

const (
	description = "Msg sends an example interactive messages."
	pluginName  = "msg"
)

// version is set via ldflags by GoReleaser.
var version = "dev"

// MsgExecutor implements the Botkube executor plugin interface.
type MsgExecutor struct{}

// Metadata returns details about the Msg plugin.
func (MsgExecutor) Metadata(context.Context) (api.MetadataOutput, error) {
	return api.MetadataOutput{
		Version:     version,
		Description: description,
	}, nil
}

// Execute returns a given command as a response.
//
//nolint:gocritic  //hugeParam: in is heavy (80 bytes); consider passing it by pointer
func (MsgExecutor) Execute(_ context.Context, in executor.ExecuteInput) (executor.ExecuteOutput, error) {
	if !in.Context.IsInteractivitySupported {
		return executor.ExecuteOutput{
			Message: api.NewCodeBlockMessage("Interactivity for this platform is not supported", true),
		}, nil
	}

	if strings.TrimSpace(in.Command) == pluginName {
		return initialMessages(), nil
	}

	msg := fmt.Sprintf("Plain command: %s", in.Command)
	return executor.ExecuteOutput{
		Message: api.NewCodeBlockMessage(msg, true),
	}, nil
}

func initialMessages() executor.ExecuteOutput {
	btnBuilder := api.NewMessageButtonBuilder()
	cmdPrefix := func(cmd string) string {
		return fmt.Sprintf("%s %s %s", api.MessageBotNamePlaceholder, pluginName, cmd)
	}

	return executor.ExecuteOutput{
		Message: api.Message{
			BaseBody: api.Body{
				Plaintext: "Showcases interactive message capabilities",
			},
			Sections: []api.Section{
				{
					Buttons: []api.Button{
						btnBuilder.ForCommandWithDescCmd("Run act1", fmt.Sprintf("%s buttons act1", pluginName)),
						btnBuilder.ForCommandWithDescCmd("Run act2", fmt.Sprintf("%s buttons act2", pluginName), api.ButtonStylePrimary),
						btnBuilder.ForCommandWithDescCmd("Run act3", fmt.Sprintf("%s buttons act3", pluginName), api.ButtonStyleDanger),
					},
				},
				{
					Buttons: []api.Button{
						btnBuilder.ForCommandWithoutDesc("Run act4", fmt.Sprintf("%s buttons act4", pluginName)),
						btnBuilder.ForCommandWithoutDesc("Run act5", fmt.Sprintf("%s buttons act5", pluginName), api.ButtonStylePrimary),
						btnBuilder.ForCommandWithoutDesc("Run act6", fmt.Sprintf("%s buttons act6", pluginName), api.ButtonStyleDanger),
					},
				},
				{
					Selects: api.Selects{
						ID: "select-id",
						Items: []api.Select{
							{
								Name:    "first",
								Command: cmdPrefix("selects first"),
								OptionGroups: []api.OptionGroup{
									{
										Name: cmdPrefix("selects first"),
										Options: []api.OptionItem{
											{Name: "BAR", Value: "BAR"},
											{Name: "BAZ", Value: "BAZ"},
											{Name: "XYZ", Value: "XYZ"},
										},
									},
								},
							},
							{
								Name:    "second",
								Command: cmdPrefix("selects second"),
								OptionGroups: []api.OptionGroup{
									{
										Name: cmdPrefix("selects second"),
										Options: []api.OptionItem{
											{Name: "BAR", Value: "BAR"},
											{Name: "BAZ", Value: "BAZ"},
											{Name: "XYZ", Value: "XYZ"},
										},
									},
									{
										Name: cmdPrefix("selects second/section2"),
										Options: []api.OptionItem{
											{Name: "123", Value: "123"},
											{Name: "456", Value: "456"},
											{Name: "789", Value: "789"},
										},
									},
								},
								// MUST be defined also under OptionGroups.Options slice.
								InitialOption: &api.OptionItem{
									Name: "789", Value: "789",
								},
							},
						},
					},
				},
			},
			PlaintextInputs: []api.LabelInput{
				{
					Command:          cmdPrefix("input-text"),
					DispatchedAction: api.DispatchInputActionOnEnter,
					Placeholder:      "String pattern to filter by",
					Text:             "Filter output",
				},
			},

			OnlyVisibleForYou: false,
			ReplaceOriginal:   false,
		},
	}
}

func (MsgExecutor) Help(context.Context) (api.Message, error) {
	msg := description
	msg += fmt.Sprintf("\nJust type `%s %s`", api.MessageBotNamePlaceholder, pluginName)

	return api.NewPlaintextMessage(msg, false), nil
}

func main() {
	executor.Serve(map[string]plugin.Plugin{
		pluginName: &executor.Plugin{
			Executor: &MsgExecutor{},
		},
	})
}
