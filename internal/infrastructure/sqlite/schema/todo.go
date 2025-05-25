package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
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
		field.String("body"),
		field.Enum("status").
			Values("NOT_STARTED", "IN_PROGRESS", "COMPLETED").
			Default("NOT_STARTED"),
	}
}

// Edges of the Todo.
func (TodoSchema) Edges() []ent.Edge {
	return nil
}
