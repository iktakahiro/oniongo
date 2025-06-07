package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Todo holds the schema definition for the Todo entity.
type TodoSchema struct {
	ent.Schema
}

func (TodoSchema) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "todo"},
	}
}

func (TodoSchema) Mixin() []ent.Mixin {
	return []ent.Mixin{
		EntityMixin{},
	}
}

// Fields of the Todo.
func (TodoSchema) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("body").Optional().Nillable(),
		field.Enum("status").
			Values("NOT_STARTED", "IN_PROGRESS", "COMPLETED").
			Default("NOT_STARTED"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("completed_at").
			Optional().
			Nillable(),
		field.Time("deleted_at").
			Optional(),
	}
}

// Edges of the Todo.
func (TodoSchema) Edges() []ent.Edge {
	return nil
}

func (TodoSchema) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted_at", "created_at"),
		index.Fields("deleted_at", "updated_at"),
		index.Fields("deleted_at", "status"),
	}
}
