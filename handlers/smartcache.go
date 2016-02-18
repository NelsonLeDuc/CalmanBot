package handlers

import (
	"database/sql"
	"fmt"
	"log"
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

	first := cached[0]
	s.monitor.ValueFor(first.ID)

	return nil
}

func (s SmartCache) CacheQuery(query, result string) int {
	queryStr := "INSERT INTO cached(query, result) VALUES($1, $2) RETURNING id"
	row := currentDB.QueryRow(queryStr, query, result)

	var id int
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
