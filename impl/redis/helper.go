package redis

import "fmt"

func rateLimitKey(key string) string {
	return fmt.Sprintf("ratelimit:%s", key)
}

func queueKey(queueName string) string {
	return fmt.Sprintf("queue:%s", queueName)
}

func processingKey(queueName string) string {
	return fmt.Sprintf("queue:%s:processing", queueName)
}

func deadLetterKey(queueName string) string {
	return fmt.Sprintf("queue:%s:dead", queueName)
}

func delayedKey(queueName string) string {
	return fmt.Sprintf("queue:%s:delayed", queueName)
}
