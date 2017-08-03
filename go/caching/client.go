package caching

import "time"

type Client interface {
	Get(cacheKey string) []byte
	Set(key string, value string, expiration time.Duration) bool
	Increment(key string) bool
}
