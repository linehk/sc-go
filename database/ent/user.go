// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/stablecog/sc-go/database/ent/user"
)

// User is the model entity for the User schema.
type User struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Email holds the value of the "email" field.
	Email string `json:"email,omitempty"`
	// StripeCustomerID holds the value of the "stripe_customer_id" field.
	StripeCustomerID string `json:"stripe_customer_id,omitempty"`
	// ActiveProductID holds the value of the "active_product_id" field.
	ActiveProductID *string `json:"active_product_id,omitempty"`
	// LastSignInAt holds the value of the "last_sign_in_at" field.
	LastSignInAt *time.Time `json:"last_sign_in_at,omitempty"`
	// LastSeenAt holds the value of the "last_seen_at" field.
	LastSeenAt time.Time `json:"last_seen_at,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserQuery when eager-loading is set.
	Edges UserEdges `json:"edges"`
}

// UserEdges holds the relations/edges for other nodes in the graph.
type UserEdges struct {
	// UserRoles holds the value of the user_roles edge.
	UserRoles []*UserRole `json:"user_roles,omitempty"`
	// Generations holds the value of the generations edge.
	Generations []*Generation `json:"generations,omitempty"`
	// Upscales holds the value of the upscales edge.
	Upscales []*Upscale `json:"upscales,omitempty"`
	// Credits holds the value of the credits edge.
	Credits []*Credit `json:"credits,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// UserRolesOrErr returns the UserRoles value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) UserRolesOrErr() ([]*UserRole, error) {
	if e.loadedTypes[0] {
		return e.UserRoles, nil
	}
	return nil, &NotLoadedError{edge: "user_roles"}
}

// GenerationsOrErr returns the Generations value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) GenerationsOrErr() ([]*Generation, error) {
	if e.loadedTypes[1] {
		return e.Generations, nil
	}
	return nil, &NotLoadedError{edge: "generations"}
}

// UpscalesOrErr returns the Upscales value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) UpscalesOrErr() ([]*Upscale, error) {
	if e.loadedTypes[2] {
		return e.Upscales, nil
	}
	return nil, &NotLoadedError{edge: "upscales"}
}

// CreditsOrErr returns the Credits value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) CreditsOrErr() ([]*Credit, error) {
	if e.loadedTypes[3] {
		return e.Credits, nil
	}
	return nil, &NotLoadedError{edge: "credits"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*User) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case user.FieldEmail, user.FieldStripeCustomerID, user.FieldActiveProductID:
			values[i] = new(sql.NullString)
		case user.FieldLastSignInAt, user.FieldLastSeenAt, user.FieldCreatedAt, user.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case user.FieldID:
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type User", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the User fields.
func (u *User) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case user.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				u.ID = *value
			}
		case user.FieldEmail:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field email", values[i])
			} else if value.Valid {
				u.Email = value.String
			}
		case user.FieldStripeCustomerID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field stripe_customer_id", values[i])
			} else if value.Valid {
				u.StripeCustomerID = value.String
			}
		case user.FieldActiveProductID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field active_product_id", values[i])
			} else if value.Valid {
				u.ActiveProductID = new(string)
				*u.ActiveProductID = value.String
			}
		case user.FieldLastSignInAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field last_sign_in_at", values[i])
			} else if value.Valid {
				u.LastSignInAt = new(time.Time)
				*u.LastSignInAt = value.Time
			}
		case user.FieldLastSeenAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field last_seen_at", values[i])
			} else if value.Valid {
				u.LastSeenAt = value.Time
			}
		case user.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				u.CreatedAt = value.Time
			}
		case user.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				u.UpdatedAt = value.Time
			}
		}
	}
	return nil
}

// QueryUserRoles queries the "user_roles" edge of the User entity.
func (u *User) QueryUserRoles() *UserRoleQuery {
	return NewUserClient(u.config).QueryUserRoles(u)
}

// QueryGenerations queries the "generations" edge of the User entity.
func (u *User) QueryGenerations() *GenerationQuery {
	return NewUserClient(u.config).QueryGenerations(u)
}

// QueryUpscales queries the "upscales" edge of the User entity.
func (u *User) QueryUpscales() *UpscaleQuery {
	return NewUserClient(u.config).QueryUpscales(u)
}

// QueryCredits queries the "credits" edge of the User entity.
func (u *User) QueryCredits() *CreditQuery {
	return NewUserClient(u.config).QueryCredits(u)
}

// Update returns a builder for updating this User.
// Note that you need to call User.Unwrap() before calling this method if this User
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *User) Update() *UserUpdateOne {
	return NewUserClient(u.config).UpdateOne(u)
}

// Unwrap unwraps the User entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (u *User) Unwrap() *User {
	_tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("ent: User is not a transactional entity")
	}
	u.config.driver = _tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *User) String() string {
	var builder strings.Builder
	builder.WriteString("User(")
	builder.WriteString(fmt.Sprintf("id=%v, ", u.ID))
	builder.WriteString("email=")
	builder.WriteString(u.Email)
	builder.WriteString(", ")
	builder.WriteString("stripe_customer_id=")
	builder.WriteString(u.StripeCustomerID)
	builder.WriteString(", ")
	if v := u.ActiveProductID; v != nil {
		builder.WriteString("active_product_id=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := u.LastSignInAt; v != nil {
		builder.WriteString("last_sign_in_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("last_seen_at=")
	builder.WriteString(u.LastSeenAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(u.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(u.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Users is a parsable slice of User.
type Users []*User

func (u Users) config(cfg config) {
	for _i := range u {
		u[_i].config = cfg
	}
}
