package app

import (
	"log/slog"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func Execute() {
	telegramApiKey := os.Getenv("TELEGRAM_API_KEY")
	bot, err := gotgbot.NewBot(telegramApiKey, nil)
	if err != nil {
		panic("Failed to load bot " + err.Error())
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			slog.Error(err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// set list of the commands for bot autocomplete
	bot.SetMyCommands([]gotgbot.BotCommand{
		{
			Command:     "start",
			Description: "Start listening for report changes",
		},
		{
			Command:     "stop",
			Description: "Stop reports changes notifications",
		},
	}, &gotgbot.SetMyCommandsOpts{})

	// set commands handlers
	dispatcher.AddHandler(handlers.NewCommand("start", Start))
	dispatcher.AddHandler(handlers.NewCommand("stop", Stop))

	slog.Info("Read subscribers, find listeners...")
	client.ReadSubscribers()
	if client.HasSubscribers() {
		StartPoller(bot)
	}

	slog.Info("Bot awaiting for subscriptions...")
	updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})

	updater.Idle()
}
