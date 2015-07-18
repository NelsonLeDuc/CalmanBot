package main

import (
    "database/sql"
    "os"
    
    _ "github.com/lib/pq"
    "fmt"
)

func DatabaseThing() string {
    dbUrl := os.Getenv("DATABASE_URL")
    database, _ := sql.Open("postgres", dbUrl)
    
//    fmt.Printf("db: %v, err: %v\n", db, err)
//    fmt.Println(os.Getenv("CLEARDB_DATABASE_URL"))
    
    var name string
    database.QueryRow("SELECT name FROM bots").Scan(&name)
    fmt.Printf("name: %v\n", name)
    
    return name
}

//
//var currentId int
//
//var todos Todos
//
//// Give us some seed data
//func init() {
//    RepoCreateTodo(Todo{Name: "Write presentation"})
//    RepoCreateTodo(Todo{Name: "Host meetup"})
//}
//
//func RepoFindTodo(id int) Todo {
//    for _, t := range todos {
//        if t.Id == id {
//            return t
//        }
//    }
//    // return empty Todo if not found
//    return Todo{}
//}
//
//func RepoCreateTodo(t Todo) Todo {
//    currentId += 1
//    t.Id = currentId
//    todos = append(todos, t)
//    return t
//}
//
//func RepoDestroyTodo(id int) error {
//    for i, t := range todos {
//        if t.Id == id {
//            todos = append(todos[:i], todos[i+1:]...)
//            return nil
//        }
//    }
//    return fmt.Errorf("Could not find Todo with id of %d to delete", id)
//}