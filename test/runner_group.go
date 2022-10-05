package test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"api-tests/store"
	"api-tests/test/validators"
	"api-tests/test/validators/rules"
)

// RunnerGroup - средство для запуска отдельных тестов в группе
type RunnerGroup struct {
	group         Group
	store         *store.Store
	wsConnections map[string]*wsConnect
}

func NewRunnerGroup(group Group) *RunnerGroup {
	return &RunnerGroup{
		group:         group,
		store:         store.NewStore(group.Init.Store),
		wsConnections: make(map[string]*wsConnect),
	}
}

// Run - Запускает выполнение отдельного теста
func (r *RunnerGroup) Run(test Case) bool {
	// Логгер с данными запроса
	logger := log.With().Str("group", r.group.Name).Str("file", test.Filename).Logger()

	if len(test.Receive.Filter) > 0 {
		msg, ok := r.receive(logger, test.Receive)
		if !ok {
			return false
		}

		// Валидация Body
		if !r.validBody(logger, test.Message, msg) {
			return false
		}
	}

	// Если есть сетевой запрос, отравляем его и проверяем ответ
	if test.Request.URL != "" {
		resp, ok := r.request(logger, test.Request)

		if !ok {
			return false
		}

		if resp != nil {
			return r.validateResponse(logger, resp, test.Response)
		}
	}

	return true
}

// validateResponse проверяет ответ
func (r *RunnerGroup) validateResponse(logger zerolog.Logger, resp *http.Response, expectResponse Response) bool {
	// вычитываем тело ответа, оно понадобится не раз
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return r.error(logger, fmt.Errorf("reading response: %w", err))
	}

	logger.Trace().
		Str("response", string(body)).
		Int("code", resp.StatusCode).
		Interface("headers", map[string][]string(resp.Header)).
		Msg("validating response")

	allValid := true

	// Валидация Headers
	if valid := r.validHeaders(logger, expectResponse.Headers, resp.Header); !valid {
		allValid = false
	}
	// Валидация HTTP Code
	if valid := r.validHTTPCode(logger, expectResponse.Code, resp.StatusCode); !valid {
		allValid = false
	}

	// Валидация Body
	if valid := r.validBody(logger, expectResponse.Body, body); !valid {
		allValid = false
	}

	if allValid {
		return true
	}

	return false
}

// validHeaders проверяет заголовки
func (r *RunnerGroup) validHeaders(logger zerolog.Logger, headerRules []rules.Rule, headers http.Header) bool {
	for index, rule := range headerRules {
		headerLogger := logger.With().
			Str("validator", fmt.Sprintf("HeaderValidator[%d]", index)).Logger()

		validator := validators.NewHeaderValidator(r.store)

		rule.Type = "string"

		if err, _ := validator.ValidHeader(rule, headers); err != nil {
			return r.error(headerLogger, fmt.Errorf("headers validate: %w", err))
		} else {
			headerLogger.Info().Msg("rule passed")
		}
	}

	return true
}

// validHTTPCode проверяет HTTP код
func (r *RunnerGroup) validHTTPCode(logger zerolog.Logger, codeRules []rules.Rule, status int) bool {
	for index, rule := range codeRules {
		httpCodeLogger := logger.With().
			Str("validator", fmt.Sprintf("HTTPCodeValidator[%d]", index)).Logger()

		validator := validators.NewHTTPCodeValidator(r.store)

		rule.Type = rules.TypeInteger

		if err := validator.ValidHTTPCode(rule, status); err != nil {
			return r.error(httpCodeLogger, fmt.Errorf("http code validate: %w", err))
		} else {
			httpCodeLogger.Info().Msg("rule passed")
		}
	}

	return true
}

// validBody проверяет тело ответа
func (r *RunnerGroup) validBody(logger zerolog.Logger, descriptions []validators.ValidatorDescr, body []byte) bool {
	for _, validatorDescr := range descriptions {
		validator, err := validators.NewBodyValidator(r.store, validatorDescr)
		if err != nil {
			return r.error(logger, fmt.Errorf("creating response validator: %w", err))
		}

		for index, rule := range validatorDescr.Rules {
			ruleLogger := logger.With().
				Str("validator", fmt.Sprintf("%s[%d]", validatorDescr.Type, index)).
				Str("rule.type", string(rule.Type)).
				Str("rule.key", rule.Key).Logger()

			if err = validator.ValidBody(rule, body); err != nil {
				return r.error(ruleLogger, fmt.Errorf("response validate: %w", err))
			} else {
				ruleLogger.Info().Msg("rule passed")
			}
		}
	}

	return true
}

// request отправляет запрос
func (r *RunnerGroup) request(logger zerolog.Logger, req Request) (*http.Response, bool) {
	switch req.Protocol {
	case "ws", "websocket":
		return r.wsRequest(logger, req)
	default:
		return r.httpRequest(logger, req)
	}
}

// httpRequest отправляет HTTP запрос
func (r *RunnerGroup) httpRequest(logger zerolog.Logger, req Request) (*http.Response, bool) {
	timeout, err := r.timeout(req.Timeout)

	if err != nil {
		return nil, r.error(logger, fmt.Errorf("preparing timeout: %w", err))
	}

	// Создаем контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request, err := r.prepareHTTPRequest(ctx, req)

	if err != nil {
		return nil, r.error(logger, fmt.Errorf("creating HTTP request: %w", err))
	}

	logger = logger.With().
		Str("method", request.Method).
		Str("url", request.URL.String()).
		Str("timeout", timeout.String()).
		Str("body", req.Body).
		Logger()

	client := &http.Client{
		// Не следовать редиректам
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// выполняем HTTP запрос
	resp, err := client.Do(request)
	if err != nil {
		return nil, r.error(logger, fmt.Errorf("sending HTTP request: %w", err))
	} else {
		logger.Trace().Msg("success HTTP request")
	}

	return resp, true
}

// prepareHTTPRequest подготавливает данные для HTTP запроса
func (r *RunnerGroup) prepareHTTPRequest(ctx context.Context, req Request) (*http.Request, error) {
	method, err := r.store.Replace(req.Method)

	if err != nil {
		return nil, fmt.Errorf("preparing method: %w", err)
	}

	url, err := r.store.Replace(req.URL)

	if err != nil {
		return nil, fmt.Errorf("preparing url: %w", err)
	}

	body, err := r.store.Replace(req.Body)

	if err != nil {
		return nil, fmt.Errorf("preparing body: %w", err)
	}

	reader := bytes.NewReader([]byte(body))
	request, err := http.NewRequestWithContext(ctx, method, url, reader)

	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	for k, v := range req.Headers {
		parsedKey, err := r.store.Replace(k)
		if err != nil {
			return nil, fmt.Errorf("parsing header key '%s': %w", k, err)
		}

		parsedValue, err := r.store.Replace(v)
		if err != nil {
			return nil, fmt.Errorf("parsing header value '%s': %w", v, err)
		}

		request.Header.Set(parsedKey, parsedValue)
	}

	return request, nil
}

// timeout парсит число из строки в time.Duration.
// Число в строке воспринимается как количество секунд.
// Принимается только целое число.
func (r *RunnerGroup) timeout(strTimeout string) (time.Duration, error) {
	timeout := 5 * time.Second

	strTimeout, err := r.store.Replace(strTimeout)

	if strTimeout == "" {
		return timeout, nil
	}

	if err != nil {
		return timeout, fmt.Errorf("preparing: %w", err)
	}

	userTimeout, err := strconv.Atoi(strTimeout)

	if err != nil {
		return timeout, fmt.Errorf("parsing integer: %w", err)
	}

	if userTimeout > 0 {
		timeout = time.Duration(userTimeout) * time.Second
	}

	return timeout, nil
}

// Flush очищает занятые ресурсы:
// 1. Закрывает websocket соединения
func (r *RunnerGroup) Flush() {
	for _, connect := range r.wsConnections {
		connect.cancel()
	}
}

func (r *RunnerGroup) error(logger zerolog.Logger, err error) bool {
	logger.Error().Err(err).Send()
	return false
}
