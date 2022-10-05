package test

import (
	"github.com/rs/zerolog/log"
)

// Runner - средство запуска всех тестов
type Runner struct {
	success int
	errors  int
}

func NewRunner() *Runner {
	return &Runner{}
}

// Run запускает выполнение группы тестов
func (r *Runner) Run(group Group) (errors int, success int) {
	log.Trace().Str("group", group.Name).Msg("run group")

	groupRunner := NewRunnerGroup(group)
	defer groupRunner.Flush()

	for _, test := range group.Tests {
		if groupRunner.Run(test) {
			r.success++
		} else {
			r.errors++
			break
		}
	}

	return r.errors, r.success
}
