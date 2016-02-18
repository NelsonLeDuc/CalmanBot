package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/kisielk/sqlstruct"
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
	row := currentDB.QueryRow("SELECT id FROM cached WHERE query=$1 AND result=$2", query, result)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return id
	}

	row = currentDB.QueryRow("INSERT INTO cached(query, result) VALUES($1, $2) RETURNING id", query, result)
	row.Scan(&id)

	return id
}

//Temp DB
var currentDB *sql.DB

func init() {
	currentDB = connect()
}

func connect() *sql.DB {
	dbUrl := os.Getenv("DATABASE_URL")
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("[x] Could not open the connection to the database. Reason: %s", err.Error())
	}
	return database
}

func cacheFetch(whereStr string, values []interface{}) ([]Cached, error) {

	queryStr := fmt.Sprintf("SELECT %s FROM cached", sqlstruct.Columns(Cached{}))

	fmt.Println(queryStr)

	rows, err := currentDB.Query(queryStr+" "+whereStr, values...)
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
