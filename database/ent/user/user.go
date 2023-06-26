// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldStripeCustomerID holds the string denoting the stripe_customer_id field in the database.
	FieldStripeCustomerID = "stripe_customer_id"
	// FieldActiveProductID holds the string denoting the active_product_id field in the database.
	FieldActiveProductID = "active_product_id"
	// FieldLastSignInAt holds the string denoting the last_sign_in_at field in the database.
	FieldLastSignInAt = "last_sign_in_at"
	// FieldLastSeenAt holds the string denoting the last_seen_at field in the database.
	FieldLastSeenAt = "last_seen_at"
	// FieldBannedAt holds the string denoting the banned_at field in the database.
	FieldBannedAt = "banned_at"
	// FieldScheduledForDeletionOn holds the string denoting the scheduled_for_deletion_on field in the database.
	FieldScheduledForDeletionOn = "scheduled_for_deletion_on"
	// FieldDataDeletedAt holds the string denoting the data_deleted_at field in the database.
	FieldDataDeletedAt = "data_deleted_at"
	// FieldWantsEmail holds the string denoting the wants_email field in the database.
	FieldWantsEmail = "wants_email"
	// FieldDiscordID holds the string denoting the discord_id field in the database.
	FieldDiscordID = "discord_id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeGenerations holds the string denoting the generations edge name in mutations.
	EdgeGenerations = "generations"
	// EdgeUpscales holds the string denoting the upscales edge name in mutations.
	EdgeUpscales = "upscales"
	// EdgeVoiceovers holds the string denoting the voiceovers edge name in mutations.
	EdgeVoiceovers = "voiceovers"
	// EdgeCredits holds the string denoting the credits edge name in mutations.
	EdgeCredits = "credits"
	// EdgeAPITokens holds the string denoting the api_tokens edge name in mutations.
	EdgeAPITokens = "api_tokens"
	// EdgeRoles holds the string denoting the roles edge name in mutations.
	EdgeRoles = "roles"
	// Table holds the table name of the user in the database.
	Table = "users"
	// GenerationsTable is the table that holds the generations relation/edge.
	GenerationsTable = "generations"
	// GenerationsInverseTable is the table name for the Generation entity.
	// It exists in this package in order to avoid circular dependency with the "generation" package.
	GenerationsInverseTable = "generations"
	// GenerationsColumn is the table column denoting the generations relation/edge.
	GenerationsColumn = "user_id"
	// UpscalesTable is the table that holds the upscales relation/edge.
	UpscalesTable = "upscales"
	// UpscalesInverseTable is the table name for the Upscale entity.
	// It exists in this package in order to avoid circular dependency with the "upscale" package.
	UpscalesInverseTable = "upscales"
	// UpscalesColumn is the table column denoting the upscales relation/edge.
	UpscalesColumn = "user_id"
	// VoiceoversTable is the table that holds the voiceovers relation/edge.
	VoiceoversTable = "voiceovers"
	// VoiceoversInverseTable is the table name for the Voiceover entity.
	// It exists in this package in order to avoid circular dependency with the "voiceover" package.
	VoiceoversInverseTable = "voiceovers"
	// VoiceoversColumn is the table column denoting the voiceovers relation/edge.
	VoiceoversColumn = "user_id"
	// CreditsTable is the table that holds the credits relation/edge.
	CreditsTable = "credits"
	// CreditsInverseTable is the table name for the Credit entity.
	// It exists in this package in order to avoid circular dependency with the "credit" package.
	CreditsInverseTable = "credits"
	// CreditsColumn is the table column denoting the credits relation/edge.
	CreditsColumn = "user_id"
	// APITokensTable is the table that holds the api_tokens relation/edge.
	APITokensTable = "api_tokens"
	// APITokensInverseTable is the table name for the ApiToken entity.
	// It exists in this package in order to avoid circular dependency with the "apitoken" package.
	APITokensInverseTable = "api_tokens"
	// APITokensColumn is the table column denoting the api_tokens relation/edge.
	APITokensColumn = "user_id"
	// RolesTable is the table that holds the roles relation/edge. The primary key declared below.
	RolesTable = "user_role_users"
	// RolesInverseTable is the table name for the Role entity.
	// It exists in this package in order to avoid circular dependency with the "role" package.
	RolesInverseTable = "roles"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldEmail,
	FieldStripeCustomerID,
	FieldActiveProductID,
	FieldLastSignInAt,
	FieldLastSeenAt,
	FieldBannedAt,
	FieldScheduledForDeletionOn,
	FieldDataDeletedAt,
	FieldWantsEmail,
	FieldDiscordID,
	FieldCreatedAt,
	FieldUpdatedAt,
}

var (
	// RolesPrimaryKey and RolesColumn2 are the table columns denoting the
	// primary key for the roles relation (M2M).
	RolesPrimaryKey = []string{"role_id", "user_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultLastSeenAt holds the default value on creation for the "last_seen_at" field.
	DefaultLastSeenAt func() time.Time
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
