package countdown

import (
	"time"
)

type Countdown struct {
	value   time.Duration // countdown time
	elapsed time.Duration // time elapsed until last pause
	start   time.Time     // start or last resume time
	period  time.Duration // update interval
	delta   time.Duration // time elapsed since last pause
	paused  bool
	expired bool
}

func New(seconds, rate int, ch chan interface{}, data interface{}) *Countdown {
	c := &Countdown{
		value:   time.Duration(seconds) * time.Second,
		elapsed: 0,
		delta:   0,
		period:  time.Duration(1000/rate) * time.Millisecond,
		paused:  true,
		expired: false,
	}

	go c.update(ch, data)

	return c
}

func (c *Countdown) Paused() bool {
	return c.paused
}

func (c *Countdown) Pause() {
	if !c.paused {
		c.paused = true
	}
}

func (c *Countdown) Resume() {
	if c.paused {
		c.paused = false
		c.start = time.Now()
		c.elapsed += c.delta
		c.delta = 0
	}
}

func (c *Countdown) Reset() {
	c.paused = true
	c.expired = false
	c.elapsed = 0
	c.delta = 0
}

func (c *Countdown) Remaining() time.Duration {
	val := c.value - (c.elapsed + c.delta)
	if val < 0 {
		val = 0
	}
	return val
}

func (c *Countdown) Expired() bool {
	return c.expired
}

func (c *Countdown) update(ch chan interface{}, data interface{}) {
	for !c.expired {
		time.Sleep(c.period)
		if !c.paused {
			c.delta = time.Since(c.start)
			if c.Remaining() <= 0 {
				c.expired = true
			}
		}
	}
	ch <- data
}
