package redis

import (
	"context"
	"encoding/json"
	"erp-service/pkg/errors"

	goredis "github.com/redis/go-redis/v9"
)

func (r *Redis) Publish(ctx context.Context, channel string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return errors.ErrInternal("failed to marshal payload").WithError(err)
	}
	return r.client.Publish(ctx, channel, data).Err()
}

func (r *Redis) PublishRaw(ctx context.Context, channel string, data []byte) error {
	return r.client.Publish(ctx, channel, data).Err()
}

func (r *Redis) Subscribe(ctx context.Context, channels ...string) *Subscription {
	pubsub := r.client.Subscribe(ctx, channels...)
	return &Subscription{
		pubsub: pubsub,
	}
}

func (r *Redis) PSubscribe(ctx context.Context, patterns ...string) *Subscription {
	pubsub := r.client.PSubscribe(ctx, patterns...)
	return &Subscription{
		pubsub: pubsub,
	}
}

func (s *Subscription) Channel() <-chan *goredis.Message {
	return s.pubsub.Channel()
}

func (s *Subscription) Receive(ctx context.Context) (*Message, error) {
	msg, err := s.pubsub.ReceiveMessage(ctx)
	if err != nil {
		return nil, errors.ErrInternal("failed to receive message").WithError(err)
	}

	return &Message{
		Channel: msg.Channel,
		Payload: json.RawMessage(msg.Payload),
	}, nil
}

func (s *Subscription) ReceiveTimeout(ctx context.Context) (any, error) {
	return s.pubsub.ReceiveTimeout(ctx, 0)
}

func (s *Subscription) Subscribe(ctx context.Context, channels ...string) error {
	return s.pubsub.Subscribe(ctx, channels...)
}

func (s *Subscription) Unsubscribe(ctx context.Context, channels ...string) error {
	return s.pubsub.Unsubscribe(ctx, channels...)
}

func (s *Subscription) Close() error {
	return s.pubsub.Close()
}

func (r *Redis) NumSub(ctx context.Context, channels ...string) (map[string]int64, error) {
	return r.client.PubSubNumSub(ctx, channels...).Result()
}

func (r *Redis) NumPat(ctx context.Context) (int64, error) {
	return r.client.PubSubNumPat(ctx).Result()
}

func (r *Redis) Channels(ctx context.Context, pattern string) ([]string, error) {
	return r.client.PubSubChannels(ctx, pattern).Result()
}

type MessageHandler func(ctx context.Context, msg *Message) error

func (s *Subscription) Listen(ctx context.Context, handler MessageHandler) error {
	defer s.Close()

	ch := s.Channel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}

			m := &Message{
				Channel: msg.Channel,
				Payload: json.RawMessage(msg.Payload),
			}

			if err := handler(ctx, m); err != nil {

				continue
			}
		}
	}
}
