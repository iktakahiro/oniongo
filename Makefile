# Directory variables
ENT_DIR := internal/infrastructure/psql

# make ent-new name=Todo
.PHONY: ent-generate
ent-generate:
	pushd internal/infrastructure/sqlite && go run -mod=mod entc.go && popd

.PHONY: ent-new
ent-new:
	pushd internal/infrastructure/sqlite && \
		go run -mod=mod entgo.io/ent/cmd/ent new \
			--target ./schema \
			$(name) && \
		popd