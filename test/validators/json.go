package validators

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/buger/jsonparser"

	"github.com/MashinaMashina/api-tests/store"
	"github.com/MashinaMashina/api-tests/test/validators/rules"
)

type JSON struct {
	StoreBase
}

// NewJSONValidator возвращает валидатор для JSON.
// Используется как один из вариантов валидации тела ответа.
func NewJSONValidator(store *store.Store) *JSON {
	return &JSON{
		StoreBase{store: store},
	}
}

func (j *JSON) ValidBody(rule rules.Rule, body []byte) error {
	var err error
	rule, err = j.prepareRule(rule)

	if err != nil {
		return err
	}

	value, err := j.getValue(rule.Type, rule.Key, body)

	if err != nil {
		return fmt.Errorf("getting value of '%s': %w", rule.Key, err)
	}

	if err = rule.Valid(value); err != nil {
		return fmt.Errorf("field '%s': %w", rule.Key, err)
	}

	if rule.Store != nil {
		bytes, _, _, err := jsonparser.Get(body, rule.Key)

		if err != nil {
			return fmt.Errorf("get value for store: %w", err)
		}

		j.storeSave(rule, string(bytes))
	}

	switch rule.Type {
	case rules.TypeObject:
		for _, subRule := range rule.Fields {
			err = j.ValidBody(subRule, value.([]byte))
			if err != nil {
				return err
			}
		}
	case rules.TypeArray:
		array := value.([][]byte)
		for _, bytes := range array {
			for _, subRule := range rule.Fields {
				// В цикле проверяем правила без ключа (без ключа проверка всех элементов)
				if subRule.Key == "" {
					err = j.ValidBody(subRule, bytes)
					if err != nil {
						return err
					}
				}
			}
		}

		var index int
		for _, subRule := range rule.Fields {
			if subRule.Key != "" {
				index, err = strconv.Atoi(subRule.Key)
				if err != nil {
					return fmt.Errorf("get array key: %w", err)
				}

				if index >= len(array) {
					// Если указано, что поле не обязательное - нет ошибки
					if subRule.Required != nil && !*subRule.Required {
						continue
					} else {
						return fmt.Errorf("required, but not exists array index %d", index)
					}
				}

				// Получение по ключу произвели
				subRule.Key = ""

				err = j.ValidBody(subRule, array[index])
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (j *JSON) getValue(typo rules.RuleType, key string, data []byte) (interface{}, error) {
	var (
		value interface{}
		err   error
	)

	var path []string
	if key != "" {
		path = append(path, key)
	}

	switch typo {
	case rules.TypeBoolean:
		value, err = jsonparser.GetBoolean(data, path...)
	case rules.TypeString, rules.TypeJWT, rules.TypeHEX:
		value, err = jsonparser.GetString(data, path...)
	case rules.TypeInteger:
		value, err = jsonparser.GetFloat(data, path...)
	case rules.TypeObject:
		var realType jsonparser.ValueType
		value, realType, _, err = jsonparser.Get(data, path...)

		if realType != jsonparser.Object {
			err = fmt.Errorf("value is not object")
		}
	case rules.TypeArray:
		var resultErr error
		var valueSlice [][]byte

		_, err = jsonparser.ArrayEach(data, func(bytes []byte, t jsonparser.ValueType, _ int, err error) {
			// Если была ошибка, другие элементы не обрабатываем
			if resultErr != nil {
				return
			}
			// Если есть ошибка - сохраняем её
			if err != nil {
				resultErr = err
				return
			}

			if t == jsonparser.String {
				// Какой-то баг в либе, у строк пропадают кавычки и дальше всё падает
				bytes = append([]byte(`"`), bytes...)
				bytes = append(bytes, []byte(`"`)...)
			}

			// Если ошибок нет - сохраняем значение
			valueSlice = append(valueSlice, bytes)
		}, path...)

		// если была ошибка в цикле, берем её
		if err == nil && resultErr != nil {
			err = resultErr
		}

		value = valueSlice
	default:
		err = fmt.Errorf("invalid json rule type '%s'", typo)
	}

	if err != nil {
		if errors.Is(err, jsonparser.KeyPathNotFoundError) {
			return nil, nil
		}

		return nil, err
	}

	return value, nil
}
