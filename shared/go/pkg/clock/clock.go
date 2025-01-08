package clock

import "time"

type Clock interface {
	Now() time.Time
}

type systemClock struct{}

func NewSystemClock() *systemClock {
	return &systemClock{}
}

func (c *systemClock) Now() time.Time {
	return time.Now()
}

type fixedClock struct {
	time time.Time
}

func NewFixedClock(t time.Time) *fixedClock {
	return &fixedClock{time: t}
}

func (c *fixedClock) Now() time.Time {
	return c.time
}
