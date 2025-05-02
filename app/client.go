package app

import (
	"iter"
	"maps"
	"reflect"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Client struct {
	Subscribers     map[int64]*ext.Context
	SubscribersLock sync.Mutex

	PollingStarted   bool
	PollingStateLock sync.Mutex
}

func (c *Client) Subscribe(ctx *ext.Context) {
	c.SubscribersLock.Lock()
	defer c.SubscribersLock.Unlock()

	if c.Subscribers == nil {
		c.Subscribers = map[int64]*ext.Context{}
	}

	c.Subscribers[ctx.EffectiveUser.Id] = ctx
}

func (c *Client) Unsubscribe(ctx *ext.Context) {
	c.SubscribersLock.Lock()
	defer c.SubscribersLock.Unlock()

	delete(c.Subscribers, ctx.EffectiveUser.Id)
}

func (c *Client) HasSubscriber(userId int64) bool {
	return c.Subscribers[userId] != nil
}

func (c *Client) HasSubscribers() bool {
	return reflect.ValueOf(c.Subscribers).Len() > 0
}

func (c *Client) GetAllSubscribers() iter.Seq[*ext.Context] {
	return maps.Values(c.Subscribers)
}
