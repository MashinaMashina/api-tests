package validators

import (
	"github.com/MashinaMashina/api-tests/store"
	"github.com/MashinaMashina/api-tests/test/validators/rules"
)

type String struct {
	StoreBase
}

// NewStringValidator возвращает валидатор для строки.
// Используется как один из вариантов валидации тела ответа.
func NewStringValidator(store *store.Store) *String {
	return &String{
		StoreBase{store: store},
	}
}

func (s *String) ValidBody(rule rules.Rule, body []byte) error {
	var err error
	rule, err = s.prepareRule(rule)

	if err != nil {
		return err
	}

	val := string(body)
	err = rule.Valid(val)

	if err == nil {
		s.storeSave(rule, val)
	}

	return err
}
