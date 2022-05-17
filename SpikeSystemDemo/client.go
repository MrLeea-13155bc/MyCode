package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 20000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/buy", nil)
			if err != nil {
				log.Println(i, "make err", err)
				return
			}
			req.Header.Set("id", strconv.Itoa(i))
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(i, "make err", err)
				return
			}
			ok := resp.Header.Get("ok")
			defer resp.Body.Close()
			if ok == "true" {
				log.Println(i, "is get!")
			}

		}(i)
	}
	wg.Wait()
}
