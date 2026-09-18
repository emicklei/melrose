package core

import (
	"fmt"
)

const (
	// OnceFromTo "once"
	OnceFromTo = iota
	// OnceFromToFrom "once-two-way"
	OnceFromToFrom
	// RepeatFromTo "repeat"
	RepeatFromTo
	// RepeatFromToFrom "repeat-two-way"
	RepeatFromToFrom
)

// Interval is a HasValue that has a Value between [from..to] and increments with [by].
// If the end of the interval is reached then the Value is set to [from].
// The fields of an Interval are also HasValue.
type Interval struct {
	from      HasValue
	to        HasValue
	by        HasValue
	strategy  intervalStrategy
	value     int
	direction int
}

func (i *Interval) Value() any {
	return i.value
}

// Next returns and increases its value with [by].
func (i *Interval) Next() any {
	return i.strategy.next(i)
}

// NewInterval creates new Interval. The initial Value is set to [from]. Specify the repeat strategy.
func NewInterval(from, to, by HasValue, strategy int) *Interval {
	start := Int(from)
	interval := &Interval{from: from, to: to, by: by, value: start, strategy: asIntervalStrategy(strategy)}
	if step := Int(by); step != 0 {
		if step > 0 {
			interval.direction = 1
		} else {
			interval.direction = -1
		}
	}
	return interval
}

// Storex is part of Storable.
func (i Interval) Storex() string {
	name := intervalStrategyName(i.strategy.id())
	return fmt.Sprintf("interval(%s,%s,%s,'%s')", Storex(i.from), Storex(i.to), Storex(i.by), name)
}

// Inpsect is part of Inspectable
func (i Interval) Inspect(n Inspection) {
	n.Properties["value"] = i.Value()
	n.Properties["length"] = (Int(i.to)-Int(i.from))/Int(i.by) + 1
	n.Properties["direction"] = i.direction
	n.Properties["method"] = intervalStrategyName(i.strategy.id())
}

// ParseIntervalStrategy return the non-exposed strategy based on the name. If unknown then return OnceFromTo ("once").
func ParseIntervalStrategy(s string) intervalStrategy {
	if is, ok := intervalStrategies[s]; ok {
		return is
	}
	return strategyOnceFromTo{}
}
func intervalStrategyName(i int) string {
	for name, each := range intervalStrategies {
		if each.id() == i {
			return name
		}
	}
	return "?"
}
func asIntervalStrategy(i int) intervalStrategy {
	for _, each := range intervalStrategies {
		if each.id() == i {
			return each
		}
	}
	return strategyOnceFromTo{}
}

type intervalStrategy interface {
	id() int
	next(i *Interval) any
}

var intervalStrategies = map[string]intervalStrategy{
	"once":           strategyOnceFromTo{},
	"repeat":         strategyRepeatFromTo{},
	"two-way":        strategyOnceFromToFrom{},
	"repeat-two-way": strategyRepeatFromToFrom{},
}

// this will walk the interval once from [from] to [to]
type strategyOnceFromTo struct{}

func (s strategyOnceFromTo) id() int { return OnceFromTo }
func (s strategyOnceFromTo) next(i *Interval) any {
	from := Int(i.from)
	to := Int(i.to)
	step := abs(Int(i.by))

	if i.direction >= 0 {
		if i.value >= to {
			return i.value
		}
		next := i.value + step
		if next >= to {
			i.value = to
			return i.value
		}
		i.value = next
		return i.value
	}
	if i.value <= from {
		return i.value
	}
	next := i.value - step
	if next <= from {
		i.value = from
		return i.value
	}
	i.value = next
	return i.value
}

// this will walk the interval repeatedly from [from] to [to]
type strategyRepeatFromTo struct{}

func (s strategyRepeatFromTo) id() int { return RepeatFromTo }
func (s strategyRepeatFromTo) next(i *Interval) any {
	from := Int(i.from)
	to := Int(i.to)
	step := abs(Int(i.by))
	if i.direction >= 0 {
		next := i.value + step
		if next > to {
			i.value = from
			return i.value
		}
		i.value = next
		return i.value
	}
	next := i.value - step
	if next < from {
		i.value = to
		return i.value
	}
	i.value = next
	return i.value
}

// this will walk the interval once from [from] to [to] and back to [from]
type strategyOnceFromToFrom struct{}

func (s strategyOnceFromToFrom) id() int { return OnceFromToFrom }
func (s strategyOnceFromToFrom) next(i *Interval) any {
	from := Int(i.from)
	to := Int(i.to)
	step := abs(Int(i.by))
	if i.direction > 0 {
		if i.value >= to {
			i.value = to - step
			i.direction = -1
			return i.value
		}
		next := i.value + step
		if next >= to {
			i.value = to
			i.direction = -1
			return i.value
		}
		i.value = next
		return i.value
	}
	if i.direction < 0 {
		if i.value <= from {
			i.value = from
			i.direction = 0
			return i.value
		}
		next := i.value - step
		if next <= from {
			i.value = from
			i.direction = 0
			return i.value
		}
		i.value = next
		return i.value
	}
	return i.value
}

// this will walk the interval repeatedly from [from] to [to] and back to [from]
type strategyRepeatFromToFrom struct{}

func (s strategyRepeatFromToFrom) id() int { return RepeatFromToFrom }
func (s strategyRepeatFromToFrom) next(i *Interval) any {
	from := Int(i.from)
	to := Int(i.to)
	step := abs(Int(i.by))
	if i.direction > 0 {
		next := i.value + step
		if next > to {
			i.value = to
			i.direction = -1
			return i.value
		}
		i.value = next
		if i.value == to {
			i.direction = -1
		}
		return i.value
	}
	if i.direction < 0 {
		next := i.value - step
		if next < from {
			i.value = from
			i.direction = 1
			return i.value
		}
		i.value = next
		if i.value == from {
			i.direction = 1
		}
		return i.value
	}
	return i.value
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
