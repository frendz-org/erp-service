package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

const googleRegPrefix = "google_reg:%s"

func (r *Redis) googleRegKey(sessionID uuid.UUID) string {
	return fmt.Sprintf(googleRegPrefix, sessionID.String())
}

func (r *Redis) CreateGoogleRegistrationSession(ctx context.Context, session *entity.GoogleRegistrationSession, ttl time.Duration) error {
	key := r.googleRegKey(session.ID)

	data, err := json.Marshal(session)
	if err != nil {
		return errors.ErrInternal("failed to marshal google registration session").WithError(err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return errors.ErrInternal("failed to store google registration session").WithError(err)
	}

	return nil
}

func (r *Redis) GetGoogleRegistrationSession(ctx context.Context, sessionID uuid.UUID) (*entity.GoogleRegistrationSession, error) {
	key := r.googleRegKey(sessionID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.ErrNotFound("google registration session not found or expired")
		}
		return nil, errors.ErrInternal("failed to get google registration session").WithError(err)
	}

	var session entity.GoogleRegistrationSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, errors.ErrInternal("failed to unmarshal google registration session").WithError(err)
	}

	return &session, nil
}

func (r *Redis) GetAndDeleteGoogleRegistrationSession(ctx context.Context, sessionID uuid.UUID) (*entity.GoogleRegistrationSession, error) {
	key := r.googleRegKey(sessionID)

	data, err := r.client.GetDel(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.ErrNotFound("google registration session not found or expired")
		}
		return nil, errors.ErrInternal("failed to get and delete google registration session").WithError(err)
	}

	var session entity.GoogleRegistrationSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, errors.ErrInternal("failed to unmarshal google registration session").WithError(err)
	}

	return &session, nil
}

func (r *Redis) DeleteGoogleRegistrationSession(ctx context.Context, sessionID uuid.UUID) error {
	key := r.googleRegKey(sessionID)
	return r.client.Del(ctx, key).Err()
}
