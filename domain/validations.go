package domain

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// regex
	noSQL = regexp.MustCompile(`\b(ALTER|CREATE|DELETE|DROP|EXEC(UTE){0,1}|INSERT( +INTO){0,1}|MERGE|SELECT|UPDATE|UNION( +ALL){0,1})\b`)

	// errors

	// ErrContainsSQL ...
	ErrContainsSQL = errors.New("no SQL commands allowed to input")
)

// IsSQL ...
func IsSQL(value interface{}) error {
	s := value.(string)

	if noSQL.MatchString(strings.ToUpper(s)) {
		return ErrContainsSQL
	}

	return nil
}
