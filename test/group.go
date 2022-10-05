package test

// Group - описание набора отдельных тестов.
// Тесты отсортированы по имени файла
type Group struct {
	Name  string
	Init  Init
	Tests []Case
}

// Init - описание инициализации группы тестов
type Init struct {
	Store map[string]string `yaml:"store"`
}
