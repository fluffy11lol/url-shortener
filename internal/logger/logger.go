package loggerSetup

import (
	"log/slog"
	"os"
	"url-shortener/internal/logger/handlers/slogpretty"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

type Logger struct {
	*slog.Logger
}

func InitLogger(env string) *Logger {
	var logger *slog.Logger
	switch env {
	case local:
		//logger = slog.New(
		//	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		//)
		logger = setupPrettySlog()
	case dev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case prod:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return &Logger{logger}
}
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
func (l *Logger) ErrAttr(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
