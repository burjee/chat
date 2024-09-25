package libs

import (
	"fmt"
	"os"
	"time"
)

type Counter struct {
	f               *os.File
	counter         chan bool
	last_updated    time.Time
	count           int
	count_by_minute []int
}

func NewCounter(buffer_size int) *Counter {
	f, err := os.OpenFile("./counter.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &Counter{
		f:               f,
		counter:         make(chan bool, buffer_size),
		last_updated:    time.Now(),
		count:           0,
		count_by_minute: make([]int, 0),
	}
}

func (c *Counter) Start() {
	for range c.counter {
		now := time.Now()
		diff := now.Sub(c.last_updated)

		if diff > time.Second {
			b := []byte(fmt.Sprintf("%d\n", c.count))
			c.f.Write(b)
			c.f.Sync()

			c.count = 0
			c.last_updated = now
		}

		c.count += 1
	}
}

func (c *Counter) Increment() {
	c.counter <- true
}

func (c *Counter) Close() {
	close(c.counter)
	c.f.Close()
}
