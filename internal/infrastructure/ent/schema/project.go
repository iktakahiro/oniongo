package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
)

// ProjectSchema holds the schema definition for the ProjectSchema entity.
type ProjectSchema struct {
	ent.Schema
}

func (ProjectSchema) Mixin() []ent.Mixin {
	return []ent.Mixin{
		EntityMixin{},
	}
}

func (ProjectSchema) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "project"},
	}
}

// Fields of the ProjectSchema.
func (ProjectSchema) Fields() []ent.Field {
	return nil
}

// Edges of the ProjectSchema.
func (ProjectSchema) Edges() []ent.Edge {
	return nil
}
