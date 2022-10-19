package validators

import (
	"fmt"

	"github.com/MashinaMashina/api-tests/store"
)

// NewBodyValidator возвращает валидатор для тела ответа
func NewBodyValidator(store *store.Store, validator ValidatorDescr) (BodyValidator, error) {
	switch validator.Type {
	case "json":
		return NewJSONValidator(store), nil
	case "string":
		return NewStringValidator(store), nil
	default:
		return nil, fmt.Errorf("invalid validator type %s", validator.Type)
	}
}
