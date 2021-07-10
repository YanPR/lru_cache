// 测试本地Redis数据通过LRU Cache读取
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Ypr919/lru_cache/lrucache"
	redis "gopkg.in/redis.v4"
)

func main() {
	// Redis连接参数
	redisConf := &redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "369258",
		DB:       0, // use default DB
	}
	client, err := redisPing(redisConf)
	if err != nil {
		fmt.Println("数据库连接失败:", err)
		return
	}
	// 写入数据库
	for i := 0; i < 1000; i++ {
		redisSet(client, i)
	}

	lruCache := lrucache.New(50) //建立一个缓存
	t1 := time.Now()

	var cnt, cntSum int64
	var gotNum int = 5     // 开启的协程数量
	var testNum int = 5000 // 每个协程读取次数
	var wg sync.WaitGroup  // 声明一个WaitGroup变量

	for i := 0; i < gotNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done() // goroutinue完成后, WaitGroup的计数-1
			for i := 0; i < testNum; i++ {
				v := rand.Intn(200)
				res, flag := redisReadFromCache(client, lruCache, v)
				if s, ok := interface2Int(res); ok {
					if s != v {
						fmt.Printf("index:%v get:%v\r\n", v, s)
					}
				} else {
					fmt.Printf("interface2Int err from %v\r\n", res)
				}
				atomic.AddInt64(&cntSum, 1)
				if flag == 1 {
					atomic.AddInt64(&cnt, 1)
				}
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(t1)
	fmt.Println("Sum=", cntSum)
	fmt.Println("经Cache读取用时: ", elapsed)
	fmt.Println("Cache命中率: ", float64(cnt)/float64(gotNum*testNum))
}

// interface2Int 将空接口转换为int型
func interface2Int(val interface{}) (int, bool) {
	if mid, ok := val.(string); ok {
		if res, err := strconv.Atoi(mid); err == nil {
			return res, true
		}
	}
	return 0, false
}
