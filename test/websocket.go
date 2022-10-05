package test

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const (
	ConnOpened int32 = 1
	ConnClosed int32 = 0
)

type wsConnect struct {
	cancel   context.CancelFunc
	status   int32 // 0 - connection already closed, 1 - opened
	messages chan []byte
}

// receive ожидает получения сообщения из websocket соединения по фильтру
func (r *RunnerGroup) receive(logger zerolog.Logger, rec Receive) ([]byte, bool) {
	if rec.Channel == "" {
		return nil, r.error(logger, fmt.Errorf("empty receive channel name"))
	}

	logger = logger.With().Str("channel", rec.Channel).Logger()

	connection, ok := r.wsConnections[rec.Channel]
	if !ok {
		return nil, r.error(logger, fmt.Errorf("not found connection"))
	}

	timeout, err := r.timeout(rec.Timeout)

	if err != nil {
		return nil, r.error(logger, fmt.Errorf("preparing timeout: %w", err))
	}

	logger = logger.With().Str("timeout", timeout.String()).Logger()

	logger.Trace().Msg("receive websocket message")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, r.error(logger, fmt.Errorf("context deadline exceeded"))
		case msg, ok := <-connection.messages:
			if !ok {
				return nil, r.error(logger, fmt.Errorf("connection '%s' closed", rec.Channel))
			}

			// Отключенный логгер, чтобы фильтр не писал лог
			fakeLogger := logger.With().Logger().Level(zerolog.Disabled)
			if r.validBody(fakeLogger, rec.Filter, msg) {
				return msg, true
			}
		}
	}
}

// wsRequest создает websocket соединение
func (r *RunnerGroup) wsRequest(logger zerolog.Logger, req Request) (*http.Response, bool) {
	if req.Channel == "" {
		return nil, r.error(logger, fmt.Errorf("empty websocket channel name"))
	}

	// Закрываем, если соединение с таким именем уже было
	if channel, exists := r.wsConnections[req.Channel]; exists {
		channel.cancel()
	}

	url, err := r.store.Replace(req.URL)

	if err != nil {
		return nil, r.error(logger, fmt.Errorf("preparing url: %w", err))
	}

	logger = logger.With().Str("url", url).Logger()

	c, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, r.error(logger, fmt.Errorf("open websocket connect: %w", err))
	}

	ctx, cancel := context.WithCancel(context.Background())

	connection := &wsConnect{
		cancel:   cancel,
		messages: make(chan []byte, 256),
		status:   ConnOpened,
	}

	// read соединение
	go func() {
		defer close(connection.messages)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				// Если не было явного закрытия соединения
				if atomic.LoadInt32(&connection.status) != ConnClosed {
					logger.Error().Err(err).Msg("reading message")
				}
				return
			}

			connection.messages <- message
			logger.Info().Msgf("recv: %s", message)
		}
	}()

	// write соединение
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {

			// Закрытие соединения
			case <-ctx.Done():
				atomic.StoreInt32(&connection.status, ConnClosed)

				err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil && err.Error() != "websocket: close sent" {
					logger.Error().Err(err).Msg("closing websocket")
					return
				}

			// ping
			case t := <-ticker.C:
				err = c.WriteMessage(websocket.TextMessage, []byte(t.String()))
				if err != nil && atomic.LoadInt32(&connection.status) != ConnClosed {
					logger.Error().Err(err).Msg("write websocket ping")
					return
				}
			}
		}
	}()

	logger.Trace().Msg("success open websocket connection")

	r.wsConnections[req.Channel] = connection
	return resp, true
}
