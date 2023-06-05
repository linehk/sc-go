// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/stablecog/sc-go/database/ent/credittype"
)

// CreditType is the model entity for the CreditType schema.
type CreditType struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description *string `json:"description,omitempty"`
	// Amount holds the value of the "amount" field.
	Amount float32 `json:"amount,omitempty"`
	// StripeProductID holds the value of the "stripe_product_id" field.
	StripeProductID *string `json:"stripe_product_id,omitempty"`
	// Type holds the value of the "type" field.
	Type credittype.Type `json:"type,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CreditTypeQuery when eager-loading is set.
	Edges CreditTypeEdges `json:"edges"`
}

// CreditTypeEdges holds the relations/edges for other nodes in the graph.
type CreditTypeEdges struct {
	// Credits holds the value of the credits edge.
	Credits []*Credit `json:"credits,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// CreditsOrErr returns the Credits value or an error if the edge
// was not loaded in eager-loading.
func (e CreditTypeEdges) CreditsOrErr() ([]*Credit, error) {
	if e.loadedTypes[0] {
		return e.Credits, nil
	}
	return nil, &NotLoadedError{edge: "credits"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*CreditType) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case credittype.FieldAmount:
			values[i] = new(sql.NullFloat64)
		case credittype.FieldName, credittype.FieldDescription, credittype.FieldStripeProductID, credittype.FieldType:
			values[i] = new(sql.NullString)
		case credittype.FieldCreatedAt, credittype.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case credittype.FieldID:
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type CreditType", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CreditType fields.
func (ct *CreditType) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case credittype.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				ct.ID = *value
			}
		case credittype.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				ct.Name = value.String
			}
		case credittype.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				ct.Description = new(string)
				*ct.Description = value.String
			}
		case credittype.FieldAmount:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field amount", values[i])
			} else if value.Valid {
				ct.Amount = float32(value.Float64)
			}
		case credittype.FieldStripeProductID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field stripe_product_id", values[i])
			} else if value.Valid {
				ct.StripeProductID = new(string)
				*ct.StripeProductID = value.String
			}
		case credittype.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				ct.Type = credittype.Type(value.String)
			}
		case credittype.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				ct.CreatedAt = value.Time
			}
		case credittype.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				ct.UpdatedAt = value.Time
			}
		}
	}
	return nil
}

// QueryCredits queries the "credits" edge of the CreditType entity.
func (ct *CreditType) QueryCredits() *CreditQuery {
	return NewCreditTypeClient(ct.config).QueryCredits(ct)
}

// Update returns a builder for updating this CreditType.
// Note that you need to call CreditType.Unwrap() before calling this method if this CreditType
// was returned from a transaction, and the transaction was committed or rolled back.
func (ct *CreditType) Update() *CreditTypeUpdateOne {
	return NewCreditTypeClient(ct.config).UpdateOne(ct)
}

// Unwrap unwraps the CreditType entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ct *CreditType) Unwrap() *CreditType {
	_tx, ok := ct.config.driver.(*txDriver)
	if !ok {
		panic("ent: CreditType is not a transactional entity")
	}
	ct.config.driver = _tx.drv
	return ct
}

// String implements the fmt.Stringer.
func (ct *CreditType) String() string {
	var builder strings.Builder
	builder.WriteString("CreditType(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ct.ID))
	builder.WriteString("name=")
	builder.WriteString(ct.Name)
	builder.WriteString(", ")
	if v := ct.Description; v != nil {
		builder.WriteString("description=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("amount=")
	builder.WriteString(fmt.Sprintf("%v", ct.Amount))
	builder.WriteString(", ")
	if v := ct.StripeProductID; v != nil {
		builder.WriteString("stripe_product_id=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(fmt.Sprintf("%v", ct.Type))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(ct.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(ct.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// CreditTypes is a parsable slice of CreditType.
type CreditTypes []*CreditType

func (ct CreditTypes) config(cfg config) {
	for _i := range ct {
		ct[_i].config = cfg
	}
}
