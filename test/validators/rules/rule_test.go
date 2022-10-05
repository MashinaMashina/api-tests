package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ValidatingCase struct {
	Name      string
	input     interface{}
	rule      Rule
	expectErr string
}

func TestValid(t *testing.T) {
	varTrue := "true"
	varStr := "any-string"
	var5 := "5"

	cases := []ValidatingCase{
		// boolean
		{
			Name:      "Проверка на true, при переданном false",
			input:     false,
			rule:      Rule{Type: "boolean", Equal: &varTrue},
			expectErr: "must be true",
		},
		{
			Name:      "Проверка на true, при переданном true",
			input:     true,
			rule:      Rule{Type: "boolean", Equal: &varTrue},
			expectErr: "",
		},
		{
			Name:      "Проверка на не true, при переданном false",
			input:     false,
			rule:      Rule{Type: "boolean", NotEqual: &varTrue},
			expectErr: "",
		},
		{
			Name:      "Проверка на не true, при переданном true",
			input:     true,
			rule:      Rule{Type: "boolean", NotEqual: &varTrue},
			expectErr: "must be not true",
		},

		// string
		{
			Name:      "Проверка строки на полное соответствие. Положительный результат",
			input:     varStr,
			rule:      Rule{Type: "string", Equal: &varStr},
			expectErr: "",
		},
		{
			Name:      "Проверка строки на полное соответствие. Ждем ошибку",
			input:     "invalid-value",
			rule:      Rule{Type: "string", Equal: &varStr},
			expectErr: "must be 'any-string'",
		},
		{
			Name:      "Проверка префикса строки.",
			input:     "any-string-value",
			rule:      Rule{Type: "string", Prefix: &varStr},
			expectErr: "",
		},
		{
			Name:      "Проверка префикса строки. Ждем ошибку",
			input:     "any-INVALID-string-value",
			rule:      Rule{Type: "string", Prefix: &varStr},
			expectErr: "must have prefix 'any-string'",
		},
		{
			Name:      "Проверка суффикса строки.",
			input:     "value of any-string",
			rule:      Rule{Type: "string", Suffix: &varStr},
			expectErr: "",
		},
		{
			Name:      "Проверка суффикса строки. Ждем ошибку",
			input:     "any-INVALID-string-value",
			rule:      Rule{Type: "string", Suffix: &varStr},
			expectErr: "must have suffix 'any-string'",
		},

		// float
		{
			Name:      "Проверка на равенство int",
			input:     5.0,
			rule:      Rule{Type: "float", Equal: &var5},
			expectErr: "",
		},
		{
			Name:      "Проверка на равенство int. Ожидаем ошибку",
			input:     6.0,
			rule:      Rule{Type: "float", Equal: &var5},
			expectErr: "must be '5'",
		},
		{
			Name:      "Проверка на меньше, чем для int",
			input:     4.0,
			rule:      Rule{Type: "float", Less: &var5},
			expectErr: "",
		},
		{
			Name:      "Проверка на меньше, чем для int. Ожидаем ошибку",
			input:     5.0,
			rule:      Rule{Type: "float", Less: &var5},
			expectErr: "must less then '5'",
		},
		{
			Name:      "Проверка на больше, чем для int.",
			input:     6.0,
			rule:      Rule{Type: "float", Greater: &var5},
			expectErr: "",
		},
		{
			Name:      "Проверка на больше, чем для int. Ожидаем ошибку",
			input:     5.0,
			rule:      Rule{Type: "float", Greater: &var5},
			expectErr: "must greater then '5'",
		},

		// JWT
		{
			Name:      "Проверка JWT токена.",
			input:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			rule:      Rule{Type: "jwt"},
			expectErr: "",
		},
		{
			Name:      "Проверка JWT токена. Ждем ошибку",
			input:     "eyJhbGc2QT4fwpMeJf36POk6yJV_adQssw5c",
			rule:      Rule{Type: "jwt"},
			expectErr: "jwt token is invalid",
		},

		// hex
		{
			Name:      "Проверка hex.",
			input:     "6229a9dcb51dc400191b18ab",
			rule:      Rule{Type: "hex"},
			expectErr: "",
		},
		{
			Name:      "Проверка hex. Ждем ошибку",
			input:     "XYZ",
			rule:      Rule{Type: "hex"},
			expectErr: "encoding/hex: invalid byte: U+0058 'X'",
		},
	}

	for _, curCase := range cases {
		err := curCase.rule.Valid(curCase.input)

		// Если err == nil, то библа падает, поэтому две проверки
		if err != nil {
			assert.EqualError(t, err, curCase.expectErr, curCase.Name)
		} else if curCase.expectErr != "" {
			assert.Equal(t, curCase.expectErr, err, curCase.Name)
		}
	}
}
