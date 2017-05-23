package redis_client

import (
    "fmt"
    "github.com/go-redis/redis"
    // "log"
)

var client redis.Client

func ExampleNewClient() redis.Client{
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong, err := client.Ping().Result()
    if err != nil {
        fmt.Println(pong, err)
    }
    return *client
}

func GetTTL(key string) int {
    client :=  ExampleNewClient()
    val, err := client.TTL(key).Result()
    if err != nil {
        panic(err)
    }
    return int(val.Seconds())
}

func GetValue(key string) string {
    client :=  ExampleNewClient()
    val, err := client.Get(key).Result()
    if err != nil {
        panic(err)
    }
    return val
}

func ExampleClient() {
    client := ExampleNewClient()
    err := client.Set("key", "value", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := client.Get("key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)

    val2, err := client.Get("key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exists")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
    }
    // Output: key value
    // key2 does not exists
}
