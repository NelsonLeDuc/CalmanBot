package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
    
    router := NewRouter()
    
    log.Fatal(http.ListenAndServe(GetPort(), router))
}

func GetPort() string {
    port := os.Getenv("PORT")
    if port == "" {
        port = "4000"
    }
    
    return ":" + port
}
