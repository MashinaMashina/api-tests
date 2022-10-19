package validators

import (
	"fmt"
	"strings"

	"github.com/MashinaMashina/api-tests/store"
	"github.com/MashinaMashina/api-tests/test/validators/rules"
)

// Декораторы добавляют валидаторам логику взаимодействия с хранилищем,
// не усложняя код самих валидаторов.

// StoreBase - база для валидаторов с хранилищем
type StoreBase struct {
	store *store.Store
}

func (s StoreBase) storeSave(rule rules.Rule, value string) {
	if rule.Store != nil {
		s.store.Set(*rule.Store, value)
	}
}

// prepareRule заменяет переменные в правиле
func (s StoreBase) prepareRule(rule rules.Rule) (rules.Rule, error) {
	strType, err := s.store.Replace(string(rule.Type))
	if err != nil {
		return rule, fmt.Errorf("preparing rule type: %w", err)
	}

	strType = strings.ToLower(strType)
	if strType == "" {
		strType = "string"
	}

	rule.Type = rules.RuleType(strType)

	rule.Key, err = s.store.Replace(rule.Key)
	if err != nil {
		return rule, fmt.Errorf("preparing rule key: %w", err)
	}
	if rule.Equal != nil {
		*rule.Equal, err = s.store.Replace(*rule.Equal)
		if err != nil {
			return rule, fmt.Errorf("preparing rule equal: %w", err)
		}
	}
	if rule.NotEqual != nil {
		*rule.NotEqual, err = s.store.Replace(*rule.NotEqual)
		if err != nil {
			return rule, fmt.Errorf("preparing rule not-equal: %w", err)
		}
	}
	if rule.Less != nil {
		*rule.Less, err = s.store.Replace(*rule.Less)
		if err != nil {
			return rule, fmt.Errorf("preparing rule less: %w", err)
		}
	}
	if rule.Greater != nil {
		*rule.Greater, err = s.store.Replace(*rule.Greater)
		if err != nil {
			return rule, fmt.Errorf("preparing rule greater: %w", err)
		}
	}
	if rule.Prefix != nil {
		*rule.Prefix, err = s.store.Replace(*rule.Prefix)
		if err != nil {
			return rule, fmt.Errorf("preparing rule prefix: %w", err)
		}
	}
	if rule.Suffix != nil {
		*rule.Suffix, err = s.store.Replace(*rule.Suffix)
		if err != nil {
			return rule, fmt.Errorf("preparing rule suffix: %w", err)
		}
	}
	if rule.Store != nil {
		*rule.Store, err = s.store.Replace(*rule.Store)
		if err != nil {
			return rule, fmt.Errorf("preparing rule store: %w", err)
		}
	}
	if rule.Severity != nil {
		*rule.Severity, err = s.store.Replace(*rule.Severity)
		if err != nil {
			return rule, fmt.Errorf("preparing rule severity: %w", err)
		}
	}

	return rule, nil
}
