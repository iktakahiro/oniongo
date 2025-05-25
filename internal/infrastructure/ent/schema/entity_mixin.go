package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// EntityMixin implements the ent.Mixin for sharing
// UUID and created_at, updated_at field with all schemas that embed it.
type EntityMixin struct {
	mixin.Schema
}

// Fields of the EntityMixin.
func (EntityMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(func() uuid.UUID {
				id, _ := uuid.NewV7()
				return id
			}).
			Immutable().
			StorageKey("id"),
	}
}
