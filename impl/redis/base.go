package redis

import (
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}
type Message struct {
	Channel string          `json:"channel"`
	Payload json.RawMessage `json:"payload"`
}

type Subscription struct {
	pubsub *redis.PubSub
}

type RateLimitResult struct {
	Allowed   bool
	Remaining int64
	Total     int64
	ResetIn   time.Duration
}

type Job struct {
	ID        string          `json:"id"`
	Payload   json.RawMessage `json:"payload"`
	Attempts  int             `json:"attempts"`
	MaxRetry  int             `json:"max_retry"`
	CreatedAt time.Time       `json:"created_at"`
	Error     string          `json:"error,omitempty"`
}

type Semaphore struct {
	redis    *Redis
	name     string
	maxCount int64
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{client: client}
}
