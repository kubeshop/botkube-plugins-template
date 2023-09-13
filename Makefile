############
# Building #
############

build-plugins: ## Builds all plugins for all defined platforms
	goreleaser build --rm-dist --snapshot
.PHONY: build-plugins

build-plugins-single: ## Builds all plugins only for current GOOS and GOARCH.
	goreleaser build --rm-dist --single-target --snapshot
.PHONY: build-plugins-single

##############
# Generating #
##############

gen-plugin-index: ## Generate plugins YAML index file.
	go run github.com/kubeshop/botkube/hack -binaries-path "./dist" -use-archive=false
.PHONY: gen-plugin-index

###############
# Developing  #
###############

fix-lint-issues: ## Automatically fix lint issues
	go mod tidy
	go mod verify
	golangci-lint run --fix "./..."
.PHONY: fix-lint-issues

#############
# Others    #
#############

help: ## Show this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
.PHONY: help
