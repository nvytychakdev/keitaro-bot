package app

import (
	"encoding/json"
	"log/slog"
	"reflect"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Redis           *redis.Client
	Subscribers     map[int64]ext.Context
	SubscribersLock sync.Mutex

	PollingStarted   bool
	PollingStateLock sync.Mutex
}

func (c *Client) Subscribe(ctx *ext.Context) {
	c.SubscribersLock.Lock()
	defer c.SubscribersLock.Unlock()

	if c.Subscribers == nil {
		c.Subscribers = map[int64]ext.Context{}
	}

	c.Subscribers[ctx.EffectiveUser.Id] = *ctx
	c.StoreSubscribers()
}

func (c *Client) Unsubscribe(ctx *ext.Context) {
	c.SubscribersLock.Lock()
	defer c.SubscribersLock.Unlock()

	delete(c.Subscribers, ctx.EffectiveUser.Id)
	c.StoreSubscribers()
}

func (c *Client) HasSubscriber(userId int64) bool {
	_, ok := c.Subscribers[userId]
	return ok
}

func (c *Client) HasSubscribers() bool {
	return reflect.ValueOf(c.Subscribers).Len() > 0
}

func (c *Client) GetAllSubscribers() []*ext.Context {
	subs := make([]*ext.Context, 0, len(c.Subscribers))
	for _, v := range c.Subscribers {
		subs = append(subs, &v)
	}
	return subs
}

func (c *Client) StoreSubscribers() {
	data, err := json.Marshal(c.Subscribers)
	if err != nil {
		slog.Error("Not able to marshal subscribers", "Err", err.Error())
		return
	}

	err = c.Redis.Set(ctx, "subscriptions", data, 0).Err()
	if err != nil {
		slog.Error("Not able to store the value", "Err", err.Error())
		return
	}
}

func (c *Client) ReadSubscribers() {
	if c.Redis.Exists(ctx, "subscriptions").Val() == 0 {
		slog.Info("Subscriptions does not exist")
		return
	}

	data, err := c.Redis.Get(ctx, "subscriptions").Bytes()
	if err != nil {
		slog.Error("Not able to get subscriptions", "Err", err.Error())
		return
	}

	err = json.Unmarshal(data, &c.Subscribers)
	if err != nil {
		slog.Error("Not able to restore subscriptions", "Err", err.Error())
		return
	}

	slog.Info("Reading subscriptions", "Value", c.Subscribers)
}
