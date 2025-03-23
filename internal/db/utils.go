package db

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("record not found")

// StringToNullString converts a string to sql.NullString
func StringToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

// NullStringToString converts a sql.NullString to string
func NullStringToString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// TimeToNullTime converts a *time.Time to sql.NullTime
func TimeToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}

// NullTimeToTimePtr converts a sql.NullTime to *time.Time
func NullTimeToTimePtr(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// Int64ToNullInt64 converts an int64 to sql.NullInt64
func Int64ToNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

// NullInt64ToInt64 converts a sql.NullInt64 to int64 with a default value
func NullInt64ToInt64(i sql.NullInt64, defaultValue int64) int64 {
	if !i.Valid {
		return defaultValue
	}
	return i.Int64
}
