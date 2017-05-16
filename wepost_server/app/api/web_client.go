package main

import (
    "flag"
    "log"
    "net/http"
)

var addr = flag.String("addr", "localhost:8081", "http web client test")

func home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.Error(w, "Not found", 404)
        return
    }
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", 405)
        return
    }
    http.ServeFile(w, r, "home.html")
}

func main() {
    flag.Parse()
    log.SetFlags(0)
    http.HandleFunc("/", home)
    log.Fatal(http.ListenAndServe(*addr, nil))
}
