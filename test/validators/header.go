package validators

import (
	"fmt"
	"net/http"

	"api-tests/store"
	"api-tests/test/validators/rules"
)

type Header struct {
	StoreBase
}

// NewHeaderValidator возвращает валидатор для заголовков ответа
func NewHeaderValidator(store *store.Store) *Header {
	return &Header{
		StoreBase{store: store},
	}
}

func (h *Header) ValidHeader(rule rules.Rule, headers http.Header) (error, string) {
	var err error
	rule, err = h.prepareRule(rule)

	if err != nil {
		return err, ""
	}

	val := headers.Get(rule.Key)
	err = rule.Valid(val)

	if err != nil {
		return fmt.Errorf("field '%s': %w", rule.Key, err), ""
	}

	h.storeSave(rule, val)

	return nil, val
}
