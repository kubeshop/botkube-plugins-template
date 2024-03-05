package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/kubeshop/botkube/pkg/api"
	"github.com/kubeshop/botkube/pkg/api/source"
	pluginx "github.com/kubeshop/botkube/pkg/plugin"
)

var _ source.Source = (*Ticker)(nil)

// version is set via ldflags by GoReleaser.
var version = "dev"

// Config holds the source configuration.
type Config struct {
	Interval time.Duration `yaml:"interval,omitempty"`
}

// Ticker implements the Botkube executor plugin interface.
type Ticker struct {
	// specify that the source doesn't handle external requests
	source.HandleExternalRequestUnimplemented
}

// Metadata returns details about the Ticker plugin.
func (Ticker) Metadata(_ context.Context) (api.MetadataOutput, error) {
	return api.MetadataOutput{
		Version:     version,
		Description: "Emits an event at a specified interval",
	}, nil
}

// Stream sends an event after configured time duration.
//
//nolint:gocritic // hugeParam: in is heavy (120 bytes); consider passing it by pointer
func (Ticker) Stream(ctx context.Context, in source.StreamInput) (source.StreamOutput, error) {
	cfg, err := mergeConfigs(in.Configs)
	if err != nil {
		return source.StreamOutput{}, err
	}

	ticker := time.NewTicker(cfg.Interval)
	out := source.StreamOutput{
		Event: make(chan source.Event),
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
			case <-ticker.C:
				out.Event <- source.Event{
					Message: api.NewPlaintextMessage("Ticker Event", true),
				}
			}
		}
	}()

	return out, nil
}

func main() {
	source.Serve(map[string]plugin.Plugin{
		"ticker": &source.Plugin{
			Source: &Ticker{},
		},
	})
}

func mergeConfigs(configs []*source.Config) (Config, error) {
	defaults := Config{
		Interval: time.Minute,
	}

	var cfg Config
	err := pluginx.MergeSourceConfigsWithDefaults(defaults, configs, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("while parsing input configuration: %w", err)
	}
	return cfg, nil
}
