package main

import (
	"fmt"
	"strconv"

	"github.com/Ypr919/lru_cache/lrucache"
	redis "gopkg.in/redis.v4"
)

// redisPing 检测并连接Redis数据库
func redisPing(opt *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(opt)
	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := client.Ping().Result()
	return client, err
}

// redisSet 向数据库写入测试数据
func redisSet(r *redis.Client, i int) bool {
	if r == nil {
		return false
	}
	if err := r.Set("str:"+strconv.FormatInt(int64(i), 10),
		strconv.FormatInt(int64(i), 10), 0).Err(); err != nil {
		panic(err)
	}
	return true
}

// redisReadDirect 直接从数据库读取
func redisReadDirect(r *redis.Client, i int) (interface{}, bool) {
	if r == nil {
		return 0, false
	}
	value, err := r.Get("str:" + strconv.FormatInt(int64(i), 10)).Result()
	if err == redis.Nil {
		fmt.Println("str:"+strconv.FormatInt(int64(i), 10), " 不存在")
		return 0, false
	} else if err != nil {
		panic(err)
	} else {
		if res, ok := strconv.Atoi(value); ok != nil {
			return res, true
		}
		return 0, false
	}
}

// redisReadFromCache 经Cache读取,flag=0代表错误，1代表Cache返回，2代表数据库返回
func redisReadFromCache(r *redis.Client,
	c *lrucache.LRUCache, i int) (res interface{}, flag int) {
	if r == nil {
		return
	}
	// 判断Cache是否存在，有则直接返回
	if val, ok := c.Get(strconv.FormatInt(int64(i), 10)); ok > 0 {
		return val, 1
	}

	// 没有则从数据库获取并写入Cache
	value, err := r.Get("str:" + strconv.FormatInt(int64(i), 10)).Result()
	if err == redis.Nil {
		fmt.Println("str:"+strconv.FormatInt(int64(i), 10), " 不存在")
		return
	} else if err != nil {
		panic(err)
	} else {
		c.Set(strconv.FormatInt(int64(i), 10), value) // 写入Cache
		return value, 2
	}
}
