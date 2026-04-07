package create

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

func CreateLogger() *slog.Logger {
	level := slog.LevelInfo
	if os.Getenv("ENVIRONMENT") == "dev" {
		level = slog.LevelDebug
	}

	if os.Getenv("ENVIRONMENT") == "prod" {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})
		return slog.New(handler)
	}

	// dev — use devslog
	opts := &devslog.Options{
		HandlerOptions: &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
		MaxSlicePrintSize: 10,
		SortKeys:          true,
		NewLineAfterLog:   true,
		StringerFormatter: true,
		NoColor:           true,
	}
	return slog.New(devslog.NewHandler(os.Stdout, opts))
}
