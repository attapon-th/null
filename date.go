package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	// FormatDate Set default Format DateString
	FormatDate = "2006-01-02"
)

// DateString DateString string is a nullable string. It supports SQL and JSON serialization.
type DateString struct {
	sql.NullString
}

// NewDateString creates a new String
func NewDateString(s string, valid bool) DateString {
	return DateString{
		NullString: sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

// DateStringFrom creates a new String that will never be blank.
func DateStringFrom(s string) DateString {
	if t, err := time.Parse(FormatDate, s); err == nil {
		return NewDateString(t.Format(FormatDate), true)
	}
	return NewDateString(s, false)
}

// DateStringFromPtr creates a new String that be null if s is nil.
func DateStringFromPtr(s *string) DateString {
	if s == nil {
		return NewDateString("", false)
	}
	if t, err := time.Parse(FormatDate, *s); err == nil {
		return NewDateString(t.Format(FormatDate), true)
	}
	return NewDateString(*s, false)
}

func (s DateString) checkValid() bool {
	if _, err := time.Parse(FormatDate, s.String); err == nil {
		return true
	}
	return false
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (s DateString) ValueOrZero() string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null String.
func (s *DateString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		s.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &s.String); err != nil {
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	s.Valid = s.checkValid()
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this String is null.
func (s DateString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	s.String = strings.Split(s.String, "T")[0]
	return json.Marshal(s.String)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this String is null.
func (s DateString) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	s.String = strings.Split(s.String, "T")[0]
	return []byte(s.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null String if the input is a blank string.
func (s *DateString) UnmarshalText(text []byte) error {
	s.String = string(text)
	s.Valid = s.checkValid()
	return nil
}

// SetValid changes this String's value and also sets it to be non-null.
func (s *DateString) SetValid(v string) {
	s.String = v
	s.Valid = true
}

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null.
func (s DateString) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero returns true for null strings, for potential future omitempty support.
func (s DateString) IsZero() bool {
	return !s.Valid
}

// Equal returns true if both strings have the same value or are both null.
func (s DateString) Equal(other DateString) bool {
	return s.Valid == other.Valid && (!s.Valid || s.String == other.String)
}
