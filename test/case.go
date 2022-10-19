package test

import (
	"github.com/MashinaMashina/api-tests/test/validators"
	"github.com/MashinaMashina/api-tests/test/validators/rules"
)

// Case - описание отдельного теста
type Case struct {
	Filename string
	Name     string                      `yaml:"name"`
	Request  Request                     `yaml:"request"`
	Response Response                    `yaml:"response"`
	Message  []validators.ValidatorDescr `yaml:"message"` // только если Protocol==ws
	Receive  Receive                     `yaml:"receive"`
}

// Request - описание запроса
type Request struct {
	URL      string            `yaml:"url"`
	Method   string            `yaml:"method"`
	Body     string            `yaml:"body"`
	Timeout  string            `yaml:"timeout"` // string, тк может быть с переменными. В мс
	Headers  map[string]string `yaml:"headers"`
	Protocol string            `yaml:"protocol"`
	Channel  string            `yaml:"channel"` // только если Protocol==ws
}

// Response - описание валидации ответа
type Response struct {
	Headers []rules.Rule                `yaml:"headers"`
	Latency []rules.Rule                `yaml:"latency"`
	Code    []rules.Rule                `yaml:"code"`
	Body    []validators.ValidatorDescr `yaml:"body"`
}

// Receive - описание ожидаемого сообщения из websocket соединения
type Receive struct {
	Channel string                      `yaml:"channel"`
	Timeout string                      `yaml:"timeout"`
	Filter  []validators.ValidatorDescr `yaml:"filter"`
}
