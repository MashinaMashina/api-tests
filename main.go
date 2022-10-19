package main

import (
	"flag"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/MashinaMashina/api-tests/service"
)

func main() {
	noColor := flag.Bool("nocolor", false, "disable output coloring")
	loglevel := flag.String("level", "trace", "log level (panic, fatal, error, warn, info, debug, trace)")
	dir := flag.String("dir", "tests", "tests directory")
	pattern := flag.String("pattern", "", "pattern for tests")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimestampFieldName = "_t"
	zerolog.LevelFieldName = "_l"
	zerolog.MessageFieldName = "_m"
	zerolog.ErrorFieldName = "_e"
	w := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05", NoColor: *noColor}
	log.Logger = zerolog.New(w).With().Timestamp().Logger()

	level, err := zerolog.ParseLevel(strings.ToLower(*loglevel))
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to parse logging level")
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	service.Run(*dir, *pattern)
}
