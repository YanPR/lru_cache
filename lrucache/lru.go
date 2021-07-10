// LRU Cache包
package lrucache

import (
	"container/list"
	"sync"
)

// lruCache 定义了LRU缓存结构组成
type LRUCache struct {
	maxNum     int                      // Cache最大容量
	mu         sync.Mutex               // 互斥锁
	lruList    *list.List               // LRU对应双向链表
	elementMap map[string]*list.Element // 存储键-链表元素指针的map
}

type lruElement struct {
	key   string
	value interface{}
}

// New 创建一个LRU Cache
func New(maxN int) *LRUCache {
	return &LRUCache{
		maxNum:     maxN,
		mu:         *new(sync.Mutex),
		lruList:    list.New(),
		elementMap: make(map[string]*list.Element),
	}
}

// Set 向Cache写入key-value
func (c *LRUCache) Set(s string, value interface{}) int {
	if c == nil {
		return 0
	}
	// 加锁
	c.mu.Lock()
	defer c.mu.Unlock()
	// Cache存在则修改
	if v, ok := c.elementMap[s]; ok {
		c.lruList.MoveToFront(v)
		if mid, ok := v.Value.(*lruElement); ok {
			mid.value = value
			return 1
		}
		return 0
	}
	// 不存在则新建加入Cache
	newEle := c.lruList.PushFront(&lruElement{s, value})
	c.elementMap[s] = newEle
	if c.lruList.Len() > c.maxNum {
		if mid, ok := c.lruList.Remove(c.lruList.Back()).(*lruElement); ok {
			delete(c.elementMap, mid.key)
			return 2
		}
		return 0
	}
	return 2
}

// Get 从Cache获取数据
func (c *LRUCache) Get(key string) (interface{}, int) {
	if c == nil {
		return nil, 0
	}
	// 加锁
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.elementMap[key]; ok {
		c.lruList.MoveToFront(v)
		mid, ok := v.Value.(*lruElement)
		if !ok {
			return nil, 0
		}
		return mid.value, 1
	}
	return nil, 0
}

// Len 返回Cache的长度
func (c *LRUCache) Len() int {
	if c == nil {
		return 0
	}
	return c.lruList.Len()
}
