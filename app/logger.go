package app

import (
	"log/slog"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/lmittmann/tint"
)

func createTelegramLogger(ctx *ext.Context) *slog.Logger {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	logger = logger.WithGroup("Context")
	logger = logger.With("UserID", ctx.EffectiveUser.Id)
	logger = logger.With("ChatID", ctx.EffectiveChat.Id)
	logger = logger.With("Username", ctx.EffectiveChat.Username)
	return logger
}
