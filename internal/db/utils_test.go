package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringToNullString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected sql.NullString
	}{
		{
			name:     "empty string",
			input:    "",
			expected: sql.NullString{String: "", Valid: false},
		},
		{
			name:     "non-empty string",
			input:    "test",
			expected: sql.NullString{String: "test", Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToNullString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNullStringToString(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected string
	}{
		{
			name:     "invalid null string",
			input:    sql.NullString{String: "", Valid: false},
			expected: "",
		},
		{
			name:     "valid null string",
			input:    sql.NullString{String: "test", Valid: true},
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NullStringToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTimeToNullTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    *time.Time
		expected sql.NullTime
	}{
		{
			name:     "nil time",
			input:    nil,
			expected: sql.NullTime{Valid: false},
		},
		{
			name:     "valid time",
			input:    &now,
			expected: sql.NullTime{Time: now, Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimeToNullTime(tt.input)
			assert.Equal(t, tt.expected.Valid, result.Valid)
			if tt.expected.Valid {
				assert.Equal(t, tt.expected.Time.Unix(), result.Time.Unix())
			}
		})
	}
}

func TestNullTimeToTimePtr(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    sql.NullTime
		expected *time.Time
	}{
		{
			name:     "invalid null time",
			input:    sql.NullTime{Valid: false},
			expected: nil,
		},
		{
			name:     "valid null time",
			input:    sql.NullTime{Time: now, Valid: true},
			expected: &now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NullTimeToTimePtr(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expected.Unix(), result.Unix())
			}
		})
	}
}

func TestInt64ToNullInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected sql.NullInt64
	}{
		{
			name:     "zero value",
			input:    0,
			expected: sql.NullInt64{Int64: 0, Valid: true},
		},
		{
			name:     "positive value",
			input:    42,
			expected: sql.NullInt64{Int64: 42, Valid: true},
		},
		{
			name:     "negative value",
			input:    -42,
			expected: sql.NullInt64{Int64: -42, Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int64ToNullInt64(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNullInt64ToInt64(t *testing.T) {
	tests := []struct {
		name         string
		input        sql.NullInt64
		defaultValue int64
		expected     int64
	}{
		{
			name:         "invalid null int64",
			input:        sql.NullInt64{Valid: false},
			defaultValue: -1,
			expected:     -1,
		},
		{
			name:         "valid null int64",
			input:        sql.NullInt64{Int64: 42, Valid: true},
			defaultValue: -1,
			expected:     42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NullInt64ToInt64(tt.input, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestErrNotFound(t *testing.T) {
	assert.Equal(t, "record not found", ErrNotFound.Error())
}
