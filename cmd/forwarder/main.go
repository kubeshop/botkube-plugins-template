package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/kubeshop/botkube/pkg/api"
	"github.com/kubeshop/botkube/pkg/api/source"
)

var _ source.Source = (*Forwarder)(nil)

// version is set via ldflags by GoReleaser.
var version = "dev"

// Forwarder implements the Botkube executor plugin interface.
type Forwarder struct {
	// specify that the source doesn't handle streaming events
	source.StreamUnimplemented
}

// Metadata returns details about the Forwarder plugin.
func (Forwarder) Metadata(_ context.Context) (api.MetadataOutput, error) {
	return api.MetadataOutput{
		Version:     version,
		Description: "Emits an event every time a message is sent as an incoming webhook request",
	}, nil
}

// payload is the incoming webhook payload.
type payload struct {
	Message string `json:"message"`
}

// HandleExternalRequest handles incoming payload and returns an event based on it.
//
//nolint:gocritic // hugeParam: in is heavy (104 bytes); consider passing it by pointer
func (Forwarder) HandleExternalRequest(_ context.Context, in source.ExternalRequestInput) (source.ExternalRequestOutput, error) {
	var p payload
	err := json.Unmarshal(in.Payload, &p)
	if err != nil {
		return source.ExternalRequestOutput{}, fmt.Errorf("while unmarshaling payload: %w", err)
	}

	if p.Message == "" {
		return source.ExternalRequestOutput{}, fmt.Errorf("message cannot be empty")
	}

	msg := fmt.Sprintf("*Incoming webhook event:* %s", p.Message)
	return source.ExternalRequestOutput{
		Event: source.Event{
			Message: api.NewPlaintextMessage(msg, true),
		},
	}, nil
}
func main() {
	source.Serve(map[string]plugin.Plugin{
		"forwarder": &source.Plugin{
			Source: &Forwarder{},
		},
	})
}
