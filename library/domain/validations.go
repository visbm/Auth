package domain

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var (
	// regex
	noSQL = regexp.MustCompile(`\b(ALTER|CREATE|DELETE|DROP|EXEC(UTE){0,1}|INSERT( +INTO){0,1}|MERGE|SELECT|UPDATE|UNION( +ALL){0,1})\b`)

	// ErrContainsSQL ...
	ErrContainsSQL = errors.New("no SQL commands allowed to input")

	// ErrPassRequirements ...
	ErrPassRequirements = errors.New("password does not meet requirements")
)

// IsSQL ...
func IsSQL(value interface{}) error {
	s := value.(string)

	if noSQL.MatchString(strings.ToUpper(s)) {
		return ErrContainsSQL
	}
	return nil
}

//Password requirements
func PassReq(value interface{}) error {
	s := value.(string)
	var number, lower, upper, special bool
	for _, s := range s {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsUpper(s):
			upper = true
		case unicode.IsLower(s):
			lower = true
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special = true
		}
	}

	if !(number && lower && upper && special) {
		return ErrPassRequirements
	}

	return nil
}
