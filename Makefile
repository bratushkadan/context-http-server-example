.PHONY: build-ci
build-ci: 
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server ./cmd/server


.PHONY: build
build: 
	@CGO_ENABLED=0 go build -o server ./cmd/server

.PHONY: lint
lint:
	@go list ./... | xargs go vet
