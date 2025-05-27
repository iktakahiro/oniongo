ENT_DIR := internal/infrastructure/ent

ENT_DIR := internal/infrastructure/ent

.PHONY: install
install: 
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
	go install github.com/vektra/mockery/v3@v3.2.5
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

.PHONY: fmt
fmt:
	golangci-lint fmt

.PHONY: lint
lint:
	golangci-lint run

# make ent-new name=Todo
.PHONY: ent-generate
ent-generate:
	pushd $(ENT_DIR) && go run -mod=mod entc.go && popd

.PHONY: ent-new
ent-new:
	pushd $(ENT_DIR) && \
		go run -mod=mod entgo.io/ent/cmd/ent new \
			--target ./schema \
			$(name) && \
		popd


MIGRATION_PATH := "file://internal/infrastructure/sqlite/migrations"

# make migrate-diff name=migration_name
.PHONY: migrate-diff
migrate-diff:
	atlas migrate diff $(name) \
		--dir $(MIGRATION_PATH) \
		--to "ent://$(ENT_DIR)/schema" \
		--dev-url "sqlite://file?mode=memory&_fk=1"

.PHONY: migrate-up
migrate-up:
	atlas migrate apply \
		--dir $(MIGRATION_PATH) \
		--url "sqlite://db/dev.db?_fk=1"

.PHONY: buf-generate
buf-generate:
	buf generate

.PHONY: mockgen
mockgen:
	mockery

.PHONY: server
server: 
	go run cmd/server/main.go

.PHONY: test
test:
	go test -v ./...
