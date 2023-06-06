// Code generated by ent, DO NOT EDIT.

package voiceoveroutput

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the voiceoveroutput type in the database.
	Label = "voiceover_output"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldAudioPath holds the string denoting the audio_path field in the database.
	FieldAudioPath = "audio_path"
	// FieldIsFavorited holds the string denoting the is_favorited field in the database.
	FieldIsFavorited = "is_favorited"
	// FieldVoiceoverID holds the string denoting the voiceover_id field in the database.
	FieldVoiceoverID = "voiceover_id"
	// FieldDeletedAt holds the string denoting the deleted_at field in the database.
	FieldDeletedAt = "deleted_at"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeVoiceovers holds the string denoting the voiceovers edge name in mutations.
	EdgeVoiceovers = "voiceovers"
	// Table holds the table name of the voiceoveroutput in the database.
	Table = "voiceover_outputs"
	// VoiceoversTable is the table that holds the voiceovers relation/edge.
	VoiceoversTable = "voiceover_outputs"
	// VoiceoversInverseTable is the table name for the Voiceover entity.
	// It exists in this package in order to avoid circular dependency with the "voiceover" package.
	VoiceoversInverseTable = "voiceovers"
	// VoiceoversColumn is the table column denoting the voiceovers relation/edge.
	VoiceoversColumn = "voiceover_id"
)

// Columns holds all SQL columns for voiceoveroutput fields.
var Columns = []string{
	FieldID,
	FieldAudioPath,
	FieldIsFavorited,
	FieldVoiceoverID,
	FieldDeletedAt,
	FieldCreatedAt,
	FieldUpdatedAt,
}

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
	// DefaultIsFavorited holds the default value on creation for the "is_favorited" field.
	DefaultIsFavorited bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
