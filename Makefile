# Directory variables
ENT_DIR := internal/infrastructure/psql

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
		--dir "file://internal/infrastructure/ent/migrations" \
		--to "ent://internal/infrastructure/ent/schema" \
		--dev-url "sqlite://file?mode=memory&_fk=1"