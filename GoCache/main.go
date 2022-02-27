package main

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

func main() {
	c := cache.New(time.Second*5, time.Second*10)    // 创建一个Cache 默认过期时间20秒 清扫间隔10秒
	c.Set("Hello", "World", cache.DefaultExpiration) // cache.DefaultExpiration 使用创建Cache时的过期时间 cache.NoExpiration则该KV不会过期
	value, exists := c.Get("Hello")                  //读取数据，之后该数据仍然在cache中
	fmt.Println(value, exists)
	/*
		time.Sleep(30 * time.Second)
		value, exists = c.Get("Hello")
		fmt.Println(value, exists)
	*/ //过期的例子，OutPut： <nil> false
	c.Set("num", 100, cache.DefaultExpiration)
	// Decrement 减
	c.DecrementInt("num", 5)
	// Increment 加
	c.IncrementInt("num", 100)
	fmt.Println(c.Get("num"))
	//c.GetWithExpiration 获取数据 过期时间 是否存在
	fmt.Println(c.GetWithExpiration("Hello"))
	//返回Cache中的数量 不包含过期的
	fmt.Println(c.ItemCount())
	//获取所有可用数据
	fmt.Println(c.Items())
	//替换
	c.Replace("Hello", "Hell", cache.DefaultExpiration)
	fmt.Println(c.Get("Hello"))
	//删除所有过期数据
	c.DeleteExpired()
	//删除数据
	c.Delete("Hello")
	fmt.Println(c.Get("Hello"))
	//删除所有数据
	fmt.Println(c.Get("num"))
	c.Flush()
	fmt.Println(c.Get("num"))

}
