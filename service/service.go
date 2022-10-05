package service

import (
	"path/filepath"
	"regexp"

	"github.com/rs/zerolog/log"

	"api-tests/finder"
	"api-tests/test"
)

// Run - запускает все тесты.
// Вначале собирает информацию о всех тестах в группы,
// а после запускает группы.
func Run(dir, patternStr string) {
	groups := finder.Find(dir, "")

	if len(groups) == 0 {
		absPath, err := filepath.Abs(dir)
		if err != nil {
			log.Error().Err(err).Msgf("read tests directory '%s'", dir)
		} else {
			log.Error().Msgf("not found tests in '%s'", absPath)
		}

		return
	}

	if patternStr != "" {
		pattern, err := regexp.Compile(patternStr)
		if err != nil {
			log.Error().Err(err).Msgf("compile pattern")
			return
		}

		newGroups := make([]test.Group, 0, len(groups))
		for i := range groups {
			if pattern.MatchString(groups[i].Name) {
				newGroups = append(newGroups, groups[i])
			}
		}
		groups = newGroups
	}

	var errors, success int
	runner := test.NewRunner()
	for _, group := range groups {
		errors, success = runner.Run(group)
	}

	logger := log.With().Int("errors", errors).Int("success", success).Logger()

	if errors > 0 {
		logger.Error().Msg("tests failed")
	} else {
		logger.Info().Msg("all tests passed")
	}
}
