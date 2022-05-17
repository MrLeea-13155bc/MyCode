package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var rdb *redis.Client

func main() {
	NewRDB()
	http.HandleFunc("/buy", spike)
	http.HandleFunc("/refresh", Refresh)
	go http.ListenAndServe(":8080", nil)
	for {
		time.Sleep(time.Second / 4)
		fmt.Print("Listening   \r")
		time.Sleep(time.Second / 4)
		fmt.Print("Listening.  \r")
		time.Sleep(time.Second / 4)
		fmt.Print("Listening.. \r")
		time.Sleep(time.Second / 4)
		fmt.Print("Listening... \r")
	}
}

func spike(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("err Method"))
		return
	}
	id := req.Header.Get("id")
	ok := GetItem(id)
	resp.Header().Set("ok", strconv.FormatBool(ok))
}

func NewRDB() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Panicln("Redis Connect Default! ", err)
	}
}

func RDB() *redis.Client {
	return rdb
}

func GetItem(host string) bool {
	var ok bool
	defer RDB().Del("lock")
	for times := 0; !ok && times < 20; times++ {
		ok, _ = RDB().SetNX("lock", host, time.Millisecond*10).Result()
		if !ok {
			rand.Seed(time.Now().Unix())
			time.Sleep(time.Millisecond * time.Duration(9*rand.Intn(20)))
			continue
		}
		snum, err := RDB().Get("item").Result()
		if err != nil {
			return false
		}
		num, _ := strconv.ParseInt(snum, 10, 64)
		if num <= 0 {
			return false
		}
		_, err = RDB().Decr("item").Result()
		if err != nil {
			log.Println("buy default")
			return false
		}
	}
	return true
}

func Refresh(resp http.ResponseWriter, req *http.Request) {
	s, err := RDB().Set("item", 10, 0).Result()
	log.Println("have Refresh!", s, err)
}
