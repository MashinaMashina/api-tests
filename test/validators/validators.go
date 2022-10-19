package validators

import (
	"net/http"

	"github.com/MashinaMashina/api-tests/test/validators/rules"
)

// ValidatorDescr - описание валидатора из yaml файла
type ValidatorDescr struct {
	Type  string       `yaml:"type"`
	Rules []rules.Rule `yaml:"rules"`
}

type BodyValidator interface {
	// ValidBody валидирует body, если поле валидно,
	// то error равно nil, вторым аргументом возвращается значение поля.
	ValidBody(rule rules.Rule, body []byte) error
}

type HTTPCodeValidator interface {
	// ValidHTTPCode валидирует HTTP код ответа, если поле валидно,
	// то error равно nil, вторым аргументом возвращается значение поля.
	ValidHTTPCode(rule rules.Rule, code int) error
}

type HeaderValidator interface {
	// ValidHeader валидирует заголовки ответа, если поле валидно,
	// то error равно nil, вторым аргументом возвращается значение поля.
	ValidHeader(rule rules.Rule, headers http.Header) error
}
