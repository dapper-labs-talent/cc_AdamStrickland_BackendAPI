.PHONY: download setup test test-watch watch format lint vet shadowed audit prep compile dist build recreate create run serve

.DEFAULT_GOAL := serve

download:
	@echo "Downloading dependencies..."
	@go mod download

# install: download
# 	@echo "Install"
# 	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

setup: download

test:
	@go run github.com/onsi/ginkgo/ginkgo -r

test-watch:
	@echo
	@make test
	@echo
	@echo "---"
	@echo

watch:
	@echo "Watching..."
	@fswatch -or --event=OwnerModified cmd internal pkg test | xargs -n1 -I{} make test-watch

format:
	@echo "Formatting..."
	@go fmt ./...

lint:
	@echo "Linting..."
	# @go run github.com/golangci/golangci-lint/cmd/golangci-lint run

vet:
	@echo "Vetting..."
	@go vet ./...

shadowed:
	@echo "Analysis/Shadows..."
	@go run golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow ./...

audit: format lint vet shadowed

prep: audit

compile:
	@echo "Compiling..."
	@go build -o dapper-api ./cmd/main.go

dist:
	@echo "Dist..."
	@echo "  no-op"

build: setup test prep compile dist

bootstrap:
	@echo "Bootstrapping..."
	./dapper-api --bootstrap --reset
	@echo

recreate: build bootstrap

run:
	./dapper-api

serve: build recreate
	@echo "Serving..."
	@make run
