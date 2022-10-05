package store

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// Store - хранилище переменных.
// Используется для передачи данных от одного теста к другому и
// для конфигурирования тестов через единый файл инициализации.
type Store struct {
	data map[string]string
}

// NewStore создает объект хранилища
func NewStore(init map[string]string) *Store {
	vars := make(map[string]string)

	// Все переменные окружения, которые начинаются на TESTS_ добавляем в список
	envs := os.Environ()
	for _, env := range envs {
		if env[:6] == "TESTS_" {
			parts := strings.SplitN(env, "=", 2)

			vars[parts[0]] = ""
			if len(parts) > 1 {
				vars[parts[0]] = parts[1]
			}
		}
	}

	for k := range init {
		vars[k] = init[k]
	}

	return &Store{
		data: vars,
	}
}

// Set - устанавливает значение переменной
func (s *Store) Set(k, v string) {
	s.data[k] = v
}

// Get - получает значение переменной.
// Так же, вторым значением сообщает о наличии переменной в хранилище
func (s *Store) Get(k string) (string, bool) {
	val, ok := s.data[k]
	return val, ok
}

// Replace принимает на вход шаблон в виде строки
// и заменяет в ней переменные данными из хранилища.
// Используется формат из стандартной библиотеки text/template.
// Пример входного шаблона: 'Привет, {{.name}}!' - тут используется переменная name.
func (s *Store) Replace(pattern string) (string, error) {
	t := template.New("")
	_, err := t.Parse(pattern)

	if err != nil {
		return "", fmt.Errorf("parsing: %w", err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, s.data)

	if err != nil {
		return "", fmt.Errorf("executing: %w", err)
	}

	return buf.String(), nil
}
