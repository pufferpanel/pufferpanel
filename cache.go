package pufferpanel

import (
	"github.com/pufferpanel/pufferpanel/v3/config"
	"sync"
	"time"
)

type Message struct {
	msg  []byte
	time int64
}

type MemoryCache struct {
	Buffer   []Message
	Capacity int
	Size     int
	Lock     sync.RWMutex
}

func CreateCache() *MemoryCache {
	capacity := config.ConsoleBuffer.Value()
	if capacity <= 0 {
		capacity = 50
	}
	return &MemoryCache{
		Buffer:   make([]Message, 0),
		Capacity: capacity * 1024, //convert to KB
	}
}

func (c *MemoryCache) Read() (msg []byte, lastTime int64) {
	msg, lastTime = c.ReadFrom(0)
	return
}

func (c *MemoryCache) ReadFrom(startTime int64) (msg []byte, lastTime int64) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	result := make([]byte, 0)

	var endTime int64 = 0

	for _, v := range c.Buffer {
		if v.time > startTime {
			result = append(result, v.msg...)
			endTime = v.time
		}
	}

	if endTime == 0 {
		endTime = time.Now().Unix()
	}
	return result, endTime
}

func (c *MemoryCache) Write(b []byte) (n int, err error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	n = len(b)

	//remove data until we've gotten small enough
	var pop Message
	for c.Size+n > c.Capacity {
		pop, c.Buffer = c.Buffer[0], c.Buffer[1:]
		c.Size = c.Size - len(pop.msg)
	}

	c.Buffer = append(c.Buffer, Message{msg: b, time: time.Now().Unix()})
	c.Size = c.Size + n
	return
}
