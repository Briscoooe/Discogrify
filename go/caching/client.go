package caching

type Client interface {
	Get(cacheKey string) []byte
	Set(key string, value string) bool
	Increment(key string) bool
}
