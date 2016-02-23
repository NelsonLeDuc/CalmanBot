package groupme

import "github.com/nelsonleduc/calmanbot/config"

type GroupmeMonitor struct{}

func (g GroupmeMonitor) ValueFor(cachedID int) int {
	row := config.DB.QueryRow("SELECT sum(likes) FROM groupme_posts WHERE cache_id=$1 GROUP BY cache_id", cachedID)

	var likeCount int
	err := row.Scan(&likeCount)
	if err != nil {
		return 0
	}

	return likeCount
}
