package cache

import (
	"fmt"
	"math/rand"

	"github.com/kisielk/sqlstruct"
	"github.com/nelsonleduc/calmanbot/config"
	"github.com/nelsonleduc/calmanbot/service"
)

type Cached struct {
	ID     int    `sql:"id"`
	Query  string `sql:"query"`
	Result string `sql:"result"`
}

type SmartCache struct {
	monitor service.Monitor
}

func NewSmartCache(monitor service.Monitor) SmartCache {
	return SmartCache{monitor}
}

func (s SmartCache) CachedResponse(message string) *string {
	if s.monitor == nil {
		return nil
	}

	cached, _ := cacheFetch("WHERE query = $1", []interface{}{message})

	fmt.Print("SMART CACHE: ")
	if len(cached) == 0 {
		fmt.Println("Nothing cached")
		return nil
	}

	itemValues := make([]int, 0)
	relevantItems := make([]Cached, 0)
	sum := 0
	for _, item := range cached {
		value := s.monitor.ValueFor(item.ID)
		if value > 1 {
			sum += value
			itemValues = append(itemValues, value)
			relevantItems = append(relevantItems, item)
		}
	}

	if len(relevantItems) == 0 {
		fmt.Println("Not enough liked items")
		return nil
	} else if rand.Intn(2) == 0 {
		fmt.Println("Failed coin flip")
		return nil
	}

	index := rand.Intn(sum)
	currentIndex := 0
	selectedIndex := 0
	for idx, num := range itemValues {
		currentIndex += num
		if index < currentIndex {
			selectedIndex = idx
			break
		}
	}

	selectedItem := relevantItems[selectedIndex]

	fmt.Println(selectedItem)

	return &selectedItem.Result
}

func (s SmartCache) CacheQuery(query, result string) int {
	if s.monitor == nil {
		return 0
	}

	row := config.DB().QueryRow("SELECT id FROM cached WHERE query=$1 AND result=$2", query, result)

	var id int
	err := row.Scan(&id)
	if err == nil {
		return id
	}

	row = config.DB().QueryRow("INSERT INTO cached(query, result) VALUES($1, $2) RETURNING id", query, result)
	row.Scan(&id)

	return id
}

func (s SmartCache) LeaderboardEntries(groupID string, count int) []LeaderboardEntry {
	if s.monitor == nil {
		return []LeaderboardEntry{}
	}

	posts, err := topPosts(groupID, count)
	if err != nil {
		return []LeaderboardEntry{}
	}

	return posts
}

func topPosts(id string, limit int) ([]LeaderboardEntry, error) {
	rows, err := config.DB().Query("SELECT cached.query, groupme_posts.likes, cached.result FROM cached INNER JOIN groupme_posts ON cached.id=groupme_posts.cache_id WHERE groupme_posts.group_id = $1 ORDER BY groupme_posts.likes DESC LIMIT $2", id, limit)
	if err != nil {
		return []LeaderboardEntry{}, err
	}
	defer rows.Close()

	actions := []LeaderboardEntry{}
	for rows.Next() {
		var (
			likeCount int
			query     string
			result    string
		)
		err := rows.Scan(&query, &likeCount, &result)
		if err == nil {
			actions = append(actions, LeaderboardEntry{likeCount, query, result})
		}
	}

	return actions, nil
}

func cacheFetch(whereStr string, values []interface{}) ([]Cached, error) {

	queryStr := fmt.Sprintf("SELECT %s FROM cached", sqlstruct.Columns(Cached{}))
	rows, err := config.DB().Query(queryStr+" "+whereStr, values...)
	if err != nil {
		return []Cached{}, err
	}
	defer rows.Close()

	actions := []Cached{}
	for rows.Next() {
		var act Cached
		err := sqlstruct.Scan(&act, rows)
		if err == nil {
			actions = append(actions, act)
		}
	}

	return actions, nil
}
