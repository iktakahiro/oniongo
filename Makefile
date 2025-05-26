
# Directory variables
ENT_DIR := internal/infrastructure/psql
GOLANGCI_LINT_PATH := ~/.local/share/mise/installs/golangci-lint/2.1.6/golangci-lint-2.1.6-darwin-arm64/golangci-lint

.PHONY: install
install: 
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6


.PHONY: fmt
fmt:
	golangci-lint fmt

.PHONY: lint
lint:
	golangci-lint run

# make ent-new name=Todo
.PHONY: ent-generate
ent-generate:
	pushd internal/infrastructure/ent && go run -mod=mod entc.go && popd

.PHONY: ent-new
ent-new:
	pushd internal/infrastructure/ent && \
		go run -mod=mod entgo.io/ent/cmd/ent new \
			--target ./schema \
			$(name) && \
		popd

# make migrate-diff name=migration_name
.PHONY: migrate-diff
migrate-diff:
	atlas migrate diff $(name) \
		--dir "file://internal/infrastructure/sqlite/migrations" \
		--to "ent://internal/infrastructure/ent/schema" \
		--dev-url "sqlite://file?mode=memory&_fk=1"

