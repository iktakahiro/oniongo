// Code generated by ent, DO NOT EDIT.

package entgen

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/iktakahiro/oniongo/internal/infrastructure/ent/entgen/projectschema"
)

// ProjectSchema is the model entity for the ProjectSchema schema.
type ProjectSchema struct {
	config
	// ID of the ent.
	ID           uuid.UUID `json:"id,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ProjectSchema) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case projectschema.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ProjectSchema fields.
func (ps *ProjectSchema) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case projectschema.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				ps.ID = *value
			}
		default:
			ps.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ProjectSchema.
// This includes values selected through modifiers, order, etc.
func (ps *ProjectSchema) Value(name string) (ent.Value, error) {
	return ps.selectValues.Get(name)
}

// Update returns a builder for updating this ProjectSchema.
// Note that you need to call ProjectSchema.Unwrap() before calling this method if this ProjectSchema
// was returned from a transaction, and the transaction was committed or rolled back.
func (ps *ProjectSchema) Update() *ProjectSchemaUpdateOne {
	return NewProjectSchemaClient(ps.config).UpdateOne(ps)
}

// Unwrap unwraps the ProjectSchema entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ps *ProjectSchema) Unwrap() *ProjectSchema {
	_tx, ok := ps.config.driver.(*txDriver)
	if !ok {
		panic("entgen: ProjectSchema is not a transactional entity")
	}
	ps.config.driver = _tx.drv
	return ps
}

// String implements the fmt.Stringer.
func (ps *ProjectSchema) String() string {
	var builder strings.Builder
	builder.WriteString("ProjectSchema(")
	builder.WriteString(fmt.Sprintf("id=%v", ps.ID))
	builder.WriteByte(')')
	return builder.String()
}

// ProjectSchemas is a parsable slice of ProjectSchema.
type ProjectSchemas []*ProjectSchema
