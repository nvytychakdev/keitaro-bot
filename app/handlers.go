package app

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var ticker *time.Ticker
var tickerComplete chan bool
var client = &Client{
	Subscribers: map[int64]*ext.Context{},
}

func Start(b *gotgbot.Bot, ctx *ext.Context) error {
	client.PollingStateLock.Lock()
	defer client.PollingStateLock.Unlock()

	logger := createTelegramLogger(ctx)

	if client.HasSubscriber(ctx.EffectiveUser.Id) {
		logger.Info("Client started the bot, but already subscribed.")
		return nil
	}

	logger.Info("Bot started...")
	logger.Info("Started scanning the keitaro campaigns for user")
	b.SendMessage(ctx.EffectiveSender.ChatId, "Hello! I've started scanning your collection...", &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	if client.HasSubscribers() {
		logger.Info("Bot already has polling started", "Subs", client.Subscribers)
		client.Subscribe(ctx)
		return nil
	}

	client.Subscribe(ctx)
	if client.PollingStarted {
		logger.Error("Polling already started, skip run")
		return nil
	}

	ticker = time.NewTicker(TICKER_TIME_INTERVAL)
	tickerComplete = make(chan bool)

	go func() {

		b.SendMessage(ctx.EffectiveSender.ChatId, "Here is your reports for the day:", &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})

		client.PollingStarted = true
		trackCampaigns(b, ctx)

		logger.Info("Polling started", "Interval", TICKER_TIME_INTERVAL)

		for {
			select {
			case <-tickerComplete:
				storedActiveReports = []Report{}
				client.PollingStarted = false
				logger.Info("Polling stopped gracefully")
				return
			case <-ticker.C:
				trackCampaigns(b, ctx)
			}
		}
	}()

	return nil
}

func Stop(b *gotgbot.Bot, ctx *ext.Context) error {
	logger := createTelegramLogger(ctx)

	if !client.HasSubscriber(ctx.EffectiveUser.Id) {
		logger.Info("User attempted to stop the app when it wasn't started.", "User Id", ctx.EffectiveUser.Id)
		return nil
	}

	logger.Info("Bot stopped!", "User Id", ctx.EffectiveUser.Id, "Chat Id", ctx.EffectiveChat.Id)
	client.Unsubscribe(ctx)

	if !client.HasSubscribers() {
		ticker.Stop()
		tickerComplete <- true
	}

	b.SendMessage(ctx.EffectiveSender.ChatId, "Stopped tracking", &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}
