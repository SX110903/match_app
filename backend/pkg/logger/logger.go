package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(env string) {
	if env == "development" {
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().Timestamp().Caller().Logger()
	} else {
		log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

func Get() *zerolog.Logger {
	return &log
}

func Info() *zerolog.Event  { return log.Info() }
func Error() *zerolog.Event { return log.Error() }
func Warn() *zerolog.Event  { return log.Warn() }
func Debug() *zerolog.Event { return log.Debug() }
func Fatal() *zerolog.Event { return log.Fatal() }
