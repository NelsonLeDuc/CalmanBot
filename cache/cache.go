package cache

type QueryCache interface {
	CachedResponse(message string) *string
	CacheQuery(query, result string) int
}
