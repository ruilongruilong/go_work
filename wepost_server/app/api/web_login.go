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
    // "fmt"
    "strconv"
    "strings"

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

func manage_web_login_status(login_token string) string{
    var result string
    result = `{"status": "expired", "message": "token不合法"}`
    //key := fmt.Printf("login_token:%s", "login_token")
    s := []string{"login_token:", login_token}
    key := strings.Join(s, "")
    ttl_time := redis_client.GetTTL(key)
    if ttl_time <= 0 {
        result = `{"status": "expired", "message": "登陆信息过期"}`
        return result
    }
    user_id_str := redis_client.GetValue(key)

    user_id := 0
    var err error
    if user_id_str != "" {
        user_id, err = strconv.Atoi(user_id_str)
        if err != nil {
            user_id = 0
        }
    }
    if user_id_str == "" || user_id == 0 {
        result = `{"status": "default", "message": "等待客户端扫码"}`
    } else if user_id == -1 {
        result = `{"status": "waiting", "message": "等待客户端确认"}`
    } else if user_id == -2 {
        result = `{"status": "cancel", "message": "登陆信息已验证"}`
    }
    log.Println(result, key)
    return result
}

func web_login_status(w http.ResponseWriter, r *http.Request) {
    type Message struct {
        SessionKey  string
        Token       string
    }
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()
    for {
        log.Println("start read message")
        mt, message, err := c.ReadMessage()
        log.Println("end read message")

        var m Message
        err = json.Unmarshal(message, &m)

        rt := manage_web_login_status(m.Token)
        result := []byte(rt)

        if err != nil {
            log.Println("read:", err)
            break
        }
        log.Printf("recv: %s", message)
        err = c.WriteMessage(mt, result)
        if err != nil {
            log.Println("write:", err)
            break
        }
    }
}

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
    http.HandleFunc("/echo", echo)
    http.HandleFunc("/", home)
    http.HandleFunc("/web-login/status", web_login_status)
    log.Fatal(http.ListenAndServe(*addr, nil))
}
