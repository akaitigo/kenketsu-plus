.PHONY: build test lint format check clean generate-vapid

build:
	cd frontend && npm run build
	cd api && go build -trimpath -ldflags "-s -w" ./...

test:
	cd frontend && npx vitest run --passWithNoTests
	cd api && go test -v -race -count=1 ./...

lint:
	cd frontend && npx oxlint . && npx biome check .
	cd api && golangci-lint run ./...

format:
	cd frontend && npx biome format --write .
	cd api && gofumpt -w . && goimports -w .

check: format lint test build
	@echo "All checks passed."

clean:
	cd frontend && rm -rf dist/ .next/ coverage/ node_modules/.cache/
	cd api && go clean -cache -testcache && rm -f coverage.out

generate-vapid:
	cd api && go run ./cmd/genvapid
