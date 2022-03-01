package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	defaultTime time.Duration
	cleanTime   time.Duration
	data        map[string]Item
	mutex       sync.Mutex
}

type Item struct {
	overTime time.Duration
	info     interface{}
}

var ch = make(chan bool, 1)

//New 创建一个新的Cache
func New(defaultTime time.Duration, cleanTime time.Duration) *Cache {
	data := make(map[string]Item)
	cache := &Cache{
		cleanTime:   cleanTime,
		defaultTime: defaultTime,
		data:        data,
	}
	return cache
}

//AutoDelete 定期删除过期数据
func (c *Cache) AutoDelete() {
	ticker := time.NewTicker(c.cleanTime)
	for {
		select {
		case <-ch:
			return
		case <-ticker.C:
			c.DeleteExpired()
		}
	}
}

//DeleteExpired 清空超时的数据
func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for k, v := range c.data {
		if now > v.overTime.Nanoseconds() {
			delete(c.data, k)
		}
	}
	if len(c.data) < 1 {
		c.Stop()
	}
}

func (c *Cache) Set(k string, v interface{}, overtime time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[k] = Item{
		overTime: overtime,
		info:     v,
	}
	if len(c.data) == 1 {
		go c.AutoDelete()
	}
}

func (c *Cache) Delete(k string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, exists := c.data["k"]
	if exists {
		delete(c.data, k)
	}
	if len(c.data) < 1 {
		c.Stop()
	}
}

func (c *Cache) Stop() {
	ch <- false
}

func (c *Cache) Get(k string) interface{} {
	value, exists := c.data[k]
	if exists {
		return value.info
	}
	return nil
}

func main() {
	cache := New(time.Second*20, time.Second*10)
	fmt.Println("初始化成功")
	cache.Set("Hello", "World", time.Second*20)
	fmt.Println(cache.Get("Hello"))
	time.Sleep(time.Second * 30)
	fmt.Println(cache.Get("Hello"))
}
