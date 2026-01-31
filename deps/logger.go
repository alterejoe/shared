package create

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

func CreateLogger() *slog.Logger {

	slogOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	opts := &devslog.Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 10,
		SortKeys:          true,
		NewLineAfterLog:   true,
		StringerFormatter: true,
		NoColor:           true,
	}
	handler := devslog.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	return logger
}
