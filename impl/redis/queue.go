package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"iam-service/pkg/errors"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

func NewJob(payload any, maxRetry int) (*Job, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.ErrInternal("failed to marshal payload").WithError(err)
	}
	return &Job{
		ID:        uuid.New().String(),
		Payload:   data,
		Attempts:  0,
		MaxRetry:  maxRetry,
		CreatedAt: time.Now(),
	}, nil
}

func (r *Redis) Enqueue(ctx context.Context, queueName string, job *Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return errors.ErrInternal("failed to marshal job").WithError(err)
	}
	return r.client.RPush(ctx, queueKey(queueName), data).Err()
}

func (r *Redis) EnqueuePayload(ctx context.Context, queueName string, payload any, maxRetry int) (string, error) {
	job, err := NewJob(payload, maxRetry)
	if err != nil {
		return "", err
	}
	if err := r.Enqueue(ctx, queueName, job); err != nil {
		return "", err
	}
	return job.ID, nil
}

func (r *Redis) EnqueueDelayed(ctx context.Context, queueName string, job *Job, delay time.Duration) error {
	data, err := json.Marshal(job)
	if err != nil {
		return errors.ErrInternal("failed to marshal job").WithError(err)
	}
	score := float64(time.Now().Add(delay).Unix())
	return r.client.ZAdd(ctx, delayedKey(queueName), goredis.Z{
		Score:  score,
		Member: data,
	}).Err()
}

func (r *Redis) Dequeue(ctx context.Context, queueName string, timeout time.Duration) (*Job, error) {

	data, err := r.client.BRPopLPush(ctx, queueKey(queueName), processingKey(queueName), timeout).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.SentinelQueueEmpty
		}
		return nil, errors.ErrInternal("failed to dequeue job").WithError(err)
	}

	var job Job
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, errors.ErrInternal("failed to unmarshal job").WithError(err)
	}

	job.Attempts++
	return &job, nil
}

func (r *Redis) DequeueBlocking(ctx context.Context, queueName string) (*Job, error) {
	return r.Dequeue(ctx, queueName, 0)
}

func (r *Redis) AcknowledgeJob(ctx context.Context, queueName string, job *Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return errors.ErrInternal("failed to marshal job").WithError(err)
	}

	job.Attempts--
	origData, err := json.Marshal(job)
	if err != nil {
		return errors.ErrInternal("failed to marshal original job").WithError(err)
	}

	removed, err := r.client.LRem(ctx, processingKey(queueName), 1, origData).Result()
	if err != nil {
		return errors.ErrInternal("failed to acknowledge job").WithError(err)
	}
	if removed == 0 {

		_, err = r.client.LRem(ctx, processingKey(queueName), 1, data).Result()
		if err != nil {
			return errors.ErrInternal("failed to acknowledge job").WithError(err)
		}
	}
	return nil
}

func (r *Redis) RejectJob(ctx context.Context, queueName string, job *Job, errMsg string) error {
	job.Error = errMsg

	origJob := *job
	origJob.Attempts--
	origData, _ := json.Marshal(&origJob)
	r.client.LRem(ctx, processingKey(queueName), 1, origData)

	if job.Attempts >= job.MaxRetry {

		data, err := json.Marshal(job)
		if err != nil {
			return errors.ErrInternal("failed to marshal job").WithError(err)
		}
		return r.client.RPush(ctx, deadLetterKey(queueName), data).Err()
	}

	return r.Enqueue(ctx, queueName, job)
}

func (r *Redis) ProcessDelayed(ctx context.Context, queueName string) (int, error) {
	now := float64(time.Now().Unix())

	jobs, err := r.client.ZRangeByScore(ctx, delayedKey(queueName), &goredis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", now),
	}).Result()
	if err != nil {
		return 0, errors.ErrInternal("failed to get delayed jobs").WithError(err)
	}

	if len(jobs) == 0 {
		return 0, nil
	}

	pipe := r.client.Pipeline()
	for _, jobData := range jobs {
		pipe.RPush(ctx, queueKey(queueName), jobData)
		pipe.ZRem(ctx, delayedKey(queueName), jobData)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return 0, errors.ErrInternal("failed to move delayed jobs").WithError(err)
	}

	return len(jobs), nil
}

func (r *Redis) QueueLen(ctx context.Context, queueName string) (int64, error) {
	return r.client.LLen(ctx, queueKey(queueName)).Result()
}

func (r *Redis) ProcessingLen(ctx context.Context, queueName string) (int64, error) {
	return r.client.LLen(ctx, processingKey(queueName)).Result()
}

func (r *Redis) DeadLetterLen(ctx context.Context, queueName string) (int64, error) {
	return r.client.LLen(ctx, deadLetterKey(queueName)).Result()
}

func (r *Redis) DelayedLen(ctx context.Context, queueName string) (int64, error) {
	return r.client.ZCard(ctx, delayedKey(queueName)).Result()
}

func (r *Redis) ClearQueue(ctx context.Context, queueName string) error {
	pipe := r.client.Pipeline()
	pipe.Del(ctx, queueKey(queueName))
	pipe.Del(ctx, processingKey(queueName))
	pipe.Del(ctx, delayedKey(queueName))
	_, err := pipe.Exec(ctx)
	return err
}

func (r *Redis) RequeueProcessing(ctx context.Context, queueName string) (int64, error) {
	var count int64
	for {
		result, err := r.client.RPopLPush(ctx, processingKey(queueName), queueKey(queueName)).Result()
		if err != nil {
			if err == goredis.Nil {
				break
			}
			return count, errors.ErrInternal("failed to requeue processing job").WithError(err)
		}
		if result == "" {
			break
		}
		count++
	}
	return count, nil
}

func (r *Redis) GetDeadLetterJobs(ctx context.Context, queueName string, start, stop int64) ([]*Job, error) {
	data, err := r.client.LRange(ctx, deadLetterKey(queueName), start, stop).Result()
	if err != nil {
		return nil, errors.ErrInternal("failed to get dead letter jobs").WithError(err)
	}

	jobs := make([]*Job, 0, len(data))
	for _, d := range data {
		var job Job
		if err := json.Unmarshal([]byte(d), &job); err != nil {
			continue
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

func (r *Redis) RetryDeadLetter(ctx context.Context, queueName string, jobID string) error {
	script := goredis.NewScript(`
		local key = KEYS[1]
		local targetID = ARGV[1]
		local len = redis.call("LLEN", key)

		for i = 0, len - 1 do
			local raw = redis.call("LINDEX", key, i)
			local job = cjson.decode(raw)
			if job.id == targetID then
				redis.call("LREM", key, 1, raw)
				return raw
			end
		end
		return nil
	`)

	dlKey := deadLetterKey(queueName)
	result, err := script.Run(ctx, r.client, []string{dlKey}, jobID).Result()
	if err == goredis.Nil {
		return errors.SentinelJobNotFound
	}
	if err != nil {
		return errors.ErrInternal("failed to retry dead letter job").WithError(err)
	}

	rawStr, ok := result.(string)
	if !ok {
		return errors.SentinelJobNotFound
	}

	var job Job
	if err := json.Unmarshal([]byte(rawStr), &job); err != nil {
		return errors.ErrInternal("failed to unmarshal job").WithError(err)
	}

	job.Attempts = 0
	job.Error = ""
	return r.Enqueue(ctx, queueName, &job)
}
