package rules

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type RuleType string

const (
	TypeBoolean RuleType = "boolean"
	TypeFloat   RuleType = "float"
	TypeInteger RuleType = "integer"
	TypeString  RuleType = "string"
	TypeEmpty   RuleType = ""
	TypeJWT     RuleType = "jwt"
	TypeHEX     RuleType = "hex"
	TypeObject  RuleType = "object"
	TypeArray   RuleType = "array"
)

// Rule - правило валидации какого-либо значения
type Rule struct {
	Type     RuleType `yaml:"type"`
	Key      string   `yaml:"key"`
	Equal    *string  `yaml:"equal"`
	NotEqual *string  `yaml:"not-equal"`
	Less     *string  `yaml:"less"`
	Greater  *string  `yaml:"greater"`
	Prefix   *string  `yaml:"prefix"`
	Suffix   *string  `yaml:"suffix"`
	Store    *string  `yaml:"store"`
	Severity *string  `yaml:"severity"`
	Required *bool    `yaml:"required"`
	Fields   []Rule   `yaml:"fields"`
}

// Valid функция проверяет валидность значения согласно описанным правилам.
// Если всё верно - error равно nil
func (r Rule) Valid(value interface{}) error {
	if value == nil {
		// Если указано, что поле не обязательное - нет ошибки
		if r.Required != nil && !*r.Required {
			return nil
		} else {
			return fmt.Errorf("required not exists")
		}
	}

	switch r.Type {
	case TypeBoolean:
		if v, ok := value.(bool); ok {
			return r.validBoolean(v)
		} else {
			return fmt.Errorf("is not a boolean")
		}
	case TypeFloat, TypeInteger:
		if i, isInt := value.(int); isInt {
			return r.validFloat(float64(i))
		} else if v, isFloat := value.(float64); isFloat {
			return r.validFloat(v)
		} else {
			return fmt.Errorf("is not a int or float")
		}
	case TypeString:
		if v, ok := value.(string); ok {
			return r.validString(v)
		} else {
			return fmt.Errorf("is not a string")
		}
	case TypeJWT:
		if v, ok := value.(string); ok {
			if err := r.validString(v); err != nil {
				return err
			}

			return r.validJWT(v)
		} else {
			return fmt.Errorf("is not a string")
		}
	case TypeHEX:
		if v, ok := value.(string); ok {
			if err := r.validString(v); err != nil {
				return err
			}

			return r.validHex(v)
		} else {
			return fmt.Errorf("is not a string")
		}
	case TypeObject, TypeArray:
		return nil
	default:
		return fmt.Errorf("invalid rule type '%s'", r.Type)
	}
}

func (r Rule) validBoolean(val bool) error {
	if r.Equal != nil {
		switch *r.Equal {
		case "true":
			if !val {
				return fmt.Errorf("must be true")
			}
		case "false":
			if val {
				return fmt.Errorf("must be false")
			}
		default:
			return fmt.Errorf("invalid boolean comparison with %s", *r.Equal)
		}
	}

	if r.NotEqual != nil {
		switch *r.NotEqual {
		case "true":
			if val {
				return fmt.Errorf("must be not true")
			}
		case "false":
			if !val {
				return fmt.Errorf("must be not false")
			}
		default:
			return fmt.Errorf("invalid boolean comparison with %s", *r.NotEqual)
		}
	}

	return nil
}

func (r Rule) validFloat(val float64) error {
	if r.Equal != nil {
		floatVal, err := strconv.ParseFloat(*r.Equal, 64)

		if err != nil {
			return fmt.Errorf("parsing Equal '%s' as float: %w", *r.Equal, err)
		}

		if floatVal != val {
			return fmt.Errorf("must be '%v'", floatVal)
		}
	}

	if r.NotEqual != nil {
		floatVal, err := strconv.ParseFloat(*r.NotEqual, 64)

		if err != nil {
			return fmt.Errorf("parsing NotEqual '%s' as float: %w", *r.NotEqual, err)
		}

		if floatVal == val {
			return fmt.Errorf("must not be '%v'", floatVal)
		}
	}

	if r.Less != nil {
		floatVal, err := strconv.ParseFloat(*r.Less, 64)

		if err != nil {
			return fmt.Errorf("parsing Less '%s' as float: %w", *r.Less, err)
		}

		if floatVal <= val {
			return fmt.Errorf("must less then '%v'", floatVal)
		}
	}

	if r.Greater != nil {
		floatVal, err := strconv.ParseFloat(*r.Greater, 64)

		if err != nil {
			return fmt.Errorf("parsing Greater '%s' as float: %w", *r.Greater, err)
		}

		if floatVal >= val {
			return fmt.Errorf("must greater then '%v'", floatVal)
		}
	}

	return nil
}

func (r Rule) validString(val string) error {
	if r.Equal != nil {
		if *r.Equal != val {
			return fmt.Errorf("must be '%s'", *r.Equal)
		}
	}

	if r.NotEqual != nil {
		if *r.NotEqual == val {
			return fmt.Errorf("must not be '%s'", *r.Equal)
		}
	}

	if r.Prefix != nil {
		if !strings.HasPrefix(val, *r.Prefix) {
			return fmt.Errorf("must have prefix '%s'", *r.Prefix)
		}
	}

	if r.Suffix != nil {
		if !strings.HasSuffix(val, *r.Suffix) {
			return fmt.Errorf("must have suffix '%s'", *r.Suffix)
		}
	}

	return nil
}

func (r Rule) validJWT(token string) error {
	res, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte{}, nil
	})

	if res == nil {
		return fmt.Errorf("jwt token is invalid")
	}

	return nil
}

func (_ Rule) validHex(id string) error {
	_, err := hex.DecodeString(id)

	return err
}
