package errors

import "errors"

var (
	SentinelLockNotAcquired = errors.New("lock not acquired")
	SentinelLockNotHeld     = errors.New("lock not held")

	SentinelSemaphoreFull          = errors.New("semaphore is full")
	SentinelSemaphoreTokenNotFound = errors.New("semaphore token not found")

	SentinelQueueEmpty  = errors.New("queue is empty")
	SentinelJobNotFound = errors.New("job not found")
)
