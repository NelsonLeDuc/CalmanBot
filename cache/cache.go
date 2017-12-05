package cache

type LeaderboardEntry struct {
	LikeCount int
	Query     string
}

type QueryCache interface {
	CachedResponse(message string) *string
	CacheQuery(query, result string) int
	LeaderboardEntries(groupID string, count int) []LeaderboardEntry
}
