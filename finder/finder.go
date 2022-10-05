package finder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"api-tests/test"
)

// Find ищет все тесты в папке.
// Обходит так же вложенные директории.
func Find(dir, namePrefix string) []test.Group {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}

	group := test.Group{
		Name: namePrefix,
	}

	var groups []test.Group
	for _, file := range files {
		// Временный файл
		if file.Name()[0] == '~' {
			continue
		}

		suffix := "/" + file.Name()
		path := dir + suffix
		ext := filepath.Ext(file.Name())

		if file.IsDir() {
			groups = append(groups, Find(path, namePrefix+suffix)...)
			continue
		}

		// Тесты описываются только в yaml файлах
		if ext != ".yml" && ext != ".yaml" {
			continue
		}

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			log.Error().Str("file", path).Err(err).Msg("reading test case")
			continue
		}

		if file.Name() == "init.yaml" || file.Name() == "init.yml" {
			init, err := parseInit(bytes)
			if err != nil {
				log.Error().Str("file", path).Err(err).Msg("decoding init config")
				continue
			}

			group.Init = init
			continue
		}

		testcase, err := parseCase(bytes)
		if err != nil {
			log.Error().Str("file", path).Err(err).Msg("decoding test case")
			continue
		}

		testcase.Filename = file.Name()

		group.Tests = append(group.Tests, testcase)
	}

	if len(group.Tests) != 0 {
		// Сортируем тесты по алфавиту
		sort.Slice(group.Tests, func(i, j int) bool {
			return group.Tests[i].Filename < group.Tests[j].Filename
		})

		groups = append(groups, group)
	}

	return groups
}

// parseCase разбирает отдельный тест
func parseCase(b []byte) (test.Case, error) {
	var testcase test.Case

	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)

	if err := dec.Decode(&testcase); err != nil {
		return test.Case{}, err
	}

	if testcase.Name == "" {
		return test.Case{}, fmt.Errorf("empty test name")
	}

	if testcase.Request.URL != "" && testcase.Request.Method == "" {
		testcase.Request.Method = "GET"
	}

	return testcase, nil
}

// parseInit разбирает init файл
func parseInit(b []byte) (test.Init, error) {
	var init test.Init

	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)

	if err := dec.Decode(&init); err != nil {
		return test.Init{}, err
	}

	return init, nil
}
