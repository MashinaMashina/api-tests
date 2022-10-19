package finder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MashinaMashina/api-tests/test"
	"github.com/MashinaMashina/api-tests/test/validators"
	"github.com/MashinaMashina/api-tests/test/validators/rules"
)

type testCase struct {
	Name      string
	Input     string
	Expect    test.Case
	ExpectErr error
}

func TestSimple(t *testing.T) {
	code200 := "200"
	bearer := "Bearer"
	str5 := "5"
	warning := "warning"

	testCases := []testCase{
		{
			Name:      "Проверка на ошибку при пустом названии",
			Input:     "request:",
			ExpectErr: fmt.Errorf("empty test name"),
		},
		{
			Name:  "С простым запросом и проверкой кода ответа",
			Input: "name: Тест\nrequest:\n  url: /api/auth/login\nresponse:\n  code:\n    - equal: 200",
			Expect: test.Case{
				Name: "Тест",
				Request: test.Request{
					Method: "GET",
					URL:    "/api/auth/login",
				},
				Response: test.Response{
					Code: []rules.Rule{{
						Equal: &code200,
					}},
				},
			},
		},
		{
			Name:  "С установкой метода запроса",
			Input: "name: Тест\nrequest:\n  method: POST\n  url: /any/path",
			Expect: test.Case{
				Name: "Тест",
				Request: test.Request{
					Method: "POST",
					URL:    "/any/path",
				},
			},
		},
		{
			Name:  "С установкой заголовков",
			Input: "name: Тест\nrequest:\n  url: /any/path\n  headers:\n    Content-Type: application/json",
			Expect: test.Case{
				Name: "Тест",
				Request: test.Request{
					Method: "GET",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					URL: "/any/path",
				},
			},
		},
		{
			Name:  "С установкой тела запроса",
			Input: "name: Тест\nrequest:\n  url: /any/path\n  body: '{\"username\":\"ivan\",\"password\":\"ivan\"}'",
			Expect: test.Case{
				Name: "Тест",
				Request: test.Request{
					Method: "GET",
					URL:    "/any/path",
					Body:   "{\"username\":\"ivan\",\"password\":\"ivan\"}",
				},
			},
		},
		{
			Name:  "С установкой валидатора заголовков",
			Input: "name: Тест\nresponse:\n  headers:\n    - key: Authorization\n      prefix: Bearer",
			Expect: test.Case{
				Name: "Тест",
				Response: test.Response{
					Headers: []rules.Rule{{
						Key:    "Authorization",
						Prefix: &bearer,
					}},
				},
			},
		},
		{
			Name:  "С установкой валидатора latency",
			Input: "name: Тест\nresponse:\n  latency:\n    - less: 5\n      severity: warning",
			Expect: test.Case{
				Name: "Тест",
				Response: test.Response{
					Latency: []rules.Rule{{
						Less:     &str5,
						Severity: &warning,
					}},
				},
			},
		},
		{
			Name:  "С установкой валидатора body - string - hex",
			Input: "name: Тест\nresponse:\n  body:\n    - type: string\n      rules:\n        - type: hex",
			Expect: test.Case{
				Name: "Тест",
				Response: test.Response{
					Body: []validators.ValidatorDescr{{
						Type: "string",
						Rules: []rules.Rule{{
							Type: "hex",
						}},
					}},
				},
			},
		},
	}

	for _, curCase := range testCases {
		res, err := parseCase([]byte(curCase.Input))

		if curCase.ExpectErr != nil {
			assert.EqualError(t, err, curCase.ExpectErr.Error(), curCase.Name)
		}

		assert.Equal(t, curCase.Expect, res, curCase.Name)
	}
}
