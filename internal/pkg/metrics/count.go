package metrics

import (
	"fmt"
	"sync"
)

type CounterMetric interface {
	Metric
	Increment(delta int)
}

type counter struct {
	sync.Mutex
	name  string
	count uint64
}

// Counter creates a new counter metric (or returns an existing one if one already exists with this name and category)
func Counter(category, name string) CounterMetric {
	return newCounter(category, name, false)
}

// DebugCounter creates a new debug counter metric (or returns an existing one if one already exists with this name and category)
func DebugCounter(category, name string) CounterMetric {
	return newCounter(category, name, true)
}

func newCounter(category string, name string, debug bool) CounterMetric {
	if metric := useCategory(category, debug).findMetric(name); metric != nil {
		if c, ok := metric.(CounterMetric); ok {
			return c
		}
	}
	count := &counter{
		name: name,
	}
	useCategory(category, debug).setMetric(count)
	return count
}

func (c *counter) Name() string {
	return c.name
}

func (c *counter) Value() string {
	return fmt.Sprintf("%d", c.count)
}

// Increment safely handles both positive and negative increments to prevent integer overflow
func (c *counter) Increment(delta int) {
	c.Lock()
	defer c.Unlock()

	// Ensure delta is non-negative before casting to uint64
	if delta < 0 {
		if uint64(-delta) > c.count {
			c.count = 0 // Prevent underflow
		} else {
			c.count -= uint64(-delta)
		}
	} else {
		c.count += uint64(delta)
	}
}
