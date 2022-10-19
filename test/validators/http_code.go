package validators

import (
	"strconv"

	"github.com/MashinaMashina/api-tests/store"
	"github.com/MashinaMashina/api-tests/test/validators/rules"
)

type HTTPCode struct {
	StoreBase
}

// NewHTTPCodeValidator возвращает валидатор для HTTP кода
func NewHTTPCodeValidator(store *store.Store) *HTTPCode {
	return &HTTPCode{
		StoreBase{store: store},
	}
}

func (h *HTTPCode) ValidHTTPCode(rule rules.Rule, code int) error {
	var err error
	rule, err = h.prepareRule(rule)

	if err != nil {
		return err
	}

	h.storeSave(rule, strconv.Itoa(code))

	return rule.Valid(code)
}
