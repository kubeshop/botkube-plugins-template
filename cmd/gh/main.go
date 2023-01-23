package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/cli/go-gh"
	ghapi "github.com/cli/go-gh/pkg/api"
	"github.com/hashicorp/go-plugin"
	"github.com/kubeshop/botkube/pkg/api"
	"github.com/kubeshop/botkube/pkg/api/executor"
	"github.com/kubeshop/botkube/pkg/pluginx"
)

const (
	pluginName       = "gh"
	logsTailLines    = 150
	defaultNamespace = "default"
	gitHubAPI        = "https://api.github.com"
)

// version is set via ldflags by GoReleaser.
var version = "dev"

// Config holds the GitHub executor configuration.
type Config struct {
	GitHub struct {
		Token         string
		Repository    string
		IssueTemplate string
	}
}

// Commands defines all supported GitHub plugin commands and their flags.
type (
	Commands struct {
		Create *CreateCommand `arg:"subcommand:create"`
	}
	CreateCommand struct {
		Issue *CreateIssueCommand `arg:"subcommand:issue"`
	}
	CreateIssueCommand struct {
		Type      string `arg:"positional"`
		Namespace string `arg:"-n,--namespace"`
	}
)

// GHExecutor implements the Botkube executor plugin interface.
type GHExecutor struct{}

// Metadata returns details about the GitHub plugin.
func (*GHExecutor) Metadata(context.Context) (api.MetadataOutput, error) {
	return api.MetadataOutput{
		Version:     version,
		Description: "GitHub creates an issue on GitHub for a related Kubernetes resource.",
	}, nil
}

// Execute returns a given command as a response.
func (e *GHExecutor) Execute(ctx context.Context, in executor.ExecuteInput) (executor.ExecuteOutput, error) {
	var cfg Config
	err := pluginx.MergeExecutorConfigs(in.Configs, &cfg)
	if err != nil {
		return executor.ExecuteOutput{}, fmt.Errorf("while merging input configs: %w", err)
	}

	var cmd Commands
	err = pluginx.ParseCommand(pluginName, in.Command, &cmd)
	if err != nil {
		return executor.ExecuteOutput{}, fmt.Errorf("while parsing input command: %w", err)
	}

	if cmd.Create == nil || cmd.Create.Issue == nil {
		return executor.ExecuteOutput{
			Data: fmt.Sprintf("Usage: %s create issue KIND/NAME", pluginName),
		}, nil
	}

	issueDetails, err := getIssueDetails(ctx, cmd.Create.Issue.Namespace, cmd.Create.Issue.Type)
	if err != nil {
		return executor.ExecuteOutput{}, fmt.Errorf("while fetching logs : %w", err)
	}

	mdBody, err := renderIssueBody(cfg.GitHub.IssueTemplate, issueDetails)
	if err != nil {
		return executor.ExecuteOutput{}, fmt.Errorf("while rendering issue body: %w", err)
	}

	title := fmt.Sprintf("The `%s` malfunctions", cmd.Create.Issue.Type)
	issueURL, err := createGitHubIssue(cfg, title, mdBody)
	if err != nil {
		return executor.ExecuteOutput{}, fmt.Errorf("while creating GitHub issue: %w", err)
	}

	return executor.ExecuteOutput{
		Data: fmt.Sprintf("New issue created successfully! 🎉\n\nIssue URL: %s", issueURL),
	}, nil
}

func main() {
	executor.Serve(map[string]plugin.Plugin{
		pluginName: &executor.Plugin{
			Executor: &GHExecutor{},
		},
	})
}

func createGitHubIssue(cfg Config, title, mdBody string) (string, error) {
	client, err := gh.HTTPClient(&ghapi.ClientOptions{
		AuthToken: cfg.GitHub.Token,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/repos/%s/issues", gitHubAPI, cfg.GitHub.Repository)
	body := struct {
		Title  string   `json:"title"`
		Body   string   `json:"body"`
		Labels []string `json:"labels"`
	}{
		Title:  title,
		Body:   mdBody,
		Labels: []string{"bug"},
	}

	out, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("while marshaling request body: %w", err)
	}
	rawResp, err := client.Post(url, "application/vnd.github+json", bytes.NewReader(out))
	if err != nil {
		return "", fmt.Errorf("while creating an issue: %w", err)
	}
	defer rawResp.Body.Close()

	if rawResp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("got unexpected status code, got %v, expected %v", rawResp.StatusCode, http.StatusCreated)
	}

	resp := struct {
		URL string `json:"html_url"`
	}{}
	rawBody, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return "", fmt.Errorf("while reading body: %w", err)
	}
	err = json.Unmarshal(rawBody, &resp)
	if err != nil {
		return "", fmt.Errorf("while unmarshaling response: %w", err)
	}

	return resp.URL, nil
}

// IssueDetails holds all available information about a given issue.
type IssueDetails struct {
	Type      string
	Namespace string
	Logs      string
	Version   string
}

func getIssueDetails(ctx context.Context, namespace, name string) (IssueDetails, error) {
	if namespace == "" {
		namespace = defaultNamespace
	}
	logs, err := pluginx.ExecuteCommand(ctx, fmt.Sprintf("kubectl logs %s -n %s --tail %d", name, namespace, logsTailLines))
	if err != nil {
		return IssueDetails{}, fmt.Errorf("while getting logs: %w", err)
	}
	ver, err := pluginx.ExecuteCommand(ctx, "kubectl version -o yaml")
	if err != nil {
		return IssueDetails{}, fmt.Errorf("while getting version: %w", err)
	}

	return IssueDetails{
		Type:      name,
		Namespace: namespace,
		Logs:      logs,
		Version:   ver,
	}, nil
}

func renderIssueBody(bodyTpl string, data IssueDetails) (string, error) {
	tmpl, err := template.New("issue-body").Funcs(template.FuncMap{
		"code": func(syntax, in string) string {
			return fmt.Sprintf("\n```%s\n%s\n```\n", syntax, in)
		},
	}).Parse(bodyTpl)
	if err != nil {
		return "", fmt.Errorf("while creating template: %w", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", fmt.Errorf("while generating body: %w", err)
	}

	return body.String(), nil
}
