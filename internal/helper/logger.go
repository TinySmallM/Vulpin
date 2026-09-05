package helper

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Inspect(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		slog.Error("Ошибка сериализации", "error", err)
		return
	}
	fmt.Println(string(b))
}
