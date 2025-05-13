package app

import (
	"log/slog"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var ticker *time.Ticker
var tickerComplete chan bool

func StartPoller(b *gotgbot.Bot) error {
	if client.PollingStarted {
		slog.Error("Polling already started, skip run")
		return nil
	}

	ticker = time.NewTicker(TICKER_TIME_INTERVAL)
	tickerComplete = make(chan bool)

	go func() {

		client.PollingStarted = true
		trackCampaigns(b)

		slog.Info("Polling started", "Interval", TICKER_TIME_INTERVAL)

		for {
			select {
			case <-tickerComplete:
				storedActiveReports = []Report{}
				client.PollingStarted = false
				slog.Info("Polling stopped gracefully")
				return
			case <-ticker.C:
				trackCampaigns(b)
			}
		}
	}()

	return nil
}
