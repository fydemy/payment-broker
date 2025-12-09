.PHONY: swagger
swagger:
	swag init -g cmd/main.go --output docs --parseDependency --parseInternal

.PHONY: run
run: swagger
	go run cmd/main.go

.PHONY: build
build: swagger
	go build -o bin/app cmd/main.go

.PHONY: cli
cli:
	go run ./cmd/cli/main.go

.PHONY: build-cli
build-cli:
	go build -o bin/cli cmd/cli/main.go

.PHONY: clean
clean:
	rm -rf docs/
	rm -rf bin/