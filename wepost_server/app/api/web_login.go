// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
    "flag"
    "log"
    "net/http"
    "encoding/json"
    "strings"
    "fmt"

    "github.com/gorilla/websocket"
    "wepost_server/util"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()
    for {
        mt, message, err := c.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            break
        }
        log.Printf("recv: %s", message)
        err = c.WriteMessage(mt, message)
        if err != nil {
            log.Println("write:", err)
            break
        }
    }
}


func web_login_status(w http.ResponseWriter, r *http.Request) {
    type Message struct {
        SessionKey, Token string
    }
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()
    for {
        mt, message, err := c.ReadMessage()
        log.Println(message)
        message_str := fmt.Sprintf("%s", message)
        log.Println(strings.NewReader(message_str))

        byt := []byte(`{"Num":"b","Strs":"a"}`)
        log.Println(byt, string(byt))
        log.Println(message, string(message))

        type DataMessage struct {
            Num string
            Strs string
        }
        var d DataMessage
        d_err := json.Unmarshal(byt, &d)
        log.Println(d_err, "llllll")
        log.Println(d.Num, "dddddddd")

        var dat map[string]interface{}
        if err := json.Unmarshal(byt, &dat); err != nil {
        }
        fmt.Println(dat)
        fmt.Println(dat["SessionKey"])
        dec := json.NewDecoder(strings.NewReader(message_str))
        log.Println(dec, "dec")

        var m Message
        err = json.Unmarshal(message, &m)
        log.Println(fmt.Printf("%v: %v\n", m.SessionKey, m.Token))

        if err != nil {
            log.Println("read:", err)
            break
        }
        log.Printf("recv: %s", message)
        err = c.WriteMessage(mt, message)
        if err != nil {
            log.Println("write:", err)
            break
        }
    }
}

func home(w http.ResponseWriter, r *http.Request) {
    redis_client.ExampleClient()
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
    http.HandleFunc("/echo", echo)
    http.HandleFunc("/", home)
    http.HandleFunc("/web-login/status", web_login_status)
    log.Fatal(http.ListenAndServe(*addr, nil))
}
