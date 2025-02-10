package clock

import "time"

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func NewSystemClock() *SystemClock {
	return &SystemClock{}
}

func (c *SystemClock) Now() time.Time {
	return time.Now()
}

type FixedClock struct {
	time time.Time
}

func NewFixedClock(t time.Time) *FixedClock {
	return &FixedClock{time: t}
}

func (c *FixedClock) Now() time.Time {
	return c.time
}
