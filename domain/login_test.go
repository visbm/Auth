package domain_test

import (
	"auth/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *domain.Login
		isValid bool
	}{
		{
			name: "valid",
			u: func() *domain.Login {
				return domain.TestLogin()
			},
			isValid: true,
		}, {
			name: "empty Username",
			u: func() *domain.Login {
				login := domain.TestLogin()
				login.Username = ""
				return login
			},
			isValid: false,
		},
		{
			name: "invalid password",
			u: func() *domain.Login {
				login := domain.TestLogin()
				login.Password = "1"
				return login
			},
			isValid: false,
		},
		{
			name: "SQL Username",
			u: func() *domain.Login {
				login := domain.TestLogin()
				login.Username = "SELECT * "
				return login
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *domain.Login {
				login := domain.TestLogin()
				login.Password = ""
				return login
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *domain.Login {
				login := domain.TestLogin()
				login.Password = "1234"
				return login
			},
			isValid: false,
		},
		{
			name: "long password",
			u: func() *domain.Login {
				login := domain.TestLogin()
				login.Password = "1234567891012345678910123456789101234567891012345678910123456789101234567891012345678910123456789101234567891012345678910"
				return login
			},
			isValid: false,
		},
		{
			name: "sql password",
			u: func() *domain.Login {
				login := domain.TestLogin()
				login.Password = "ALt  9*/123#@! eR"
				return login
			},
			isValid: false,
		},

		{
			name: "wrong password",
			u: func() *domain.Login {
				login := domain.TestLogin()
				login.Password = "pass2135"
				return login
			},
			isValid: false,
		},
		{
			name: "wrong password",
			u: func() *domain.Login {
				login := domain.TestLogin()
				login.Password = "pass*213"
				return login
			},
			isValid: false,
		},

		{
			name: "wrong password",
			u: func() *domain.Login {
				login := domain.TestLogin()
				login.Password = "Pass213"
				return login
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}
