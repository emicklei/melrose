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
	from     HasValue
	to       HasValue
	by       HasValue
	strategy intervalStrategy
	value    int
}

func (i *Interval) Value() any {
	return i.value
}

// Next returns and increases its value with [by].
func (i *Interval) Next() any {
	by := Int(i.by)
	next := i.value + by
	if by < 0 {
		if next < Int(i.from) {
			i.value = Int(i.to)
			return i.value
		}
	}
	if by > 0 {
		if next > Int(i.to) {
			i.value = Int(i.from)
			return i.value
		}
	}
	i.value = next
	return i.value
}

// NewInterval creates new Interval. The initial Value is set to [from]. Specify the repeat strategy.
func NewInterval(from, to, by HasValue, strategy int) *Interval {
	start := Int(from)
	return &Interval{from: from, to: to, by: by, value: start, strategy: asIntervalStrategy(strategy)}
}

// Storex is part of Storable.
func (i Interval) Storex() string {
	if i.strategy.id() == OnceFromTo {
		return fmt.Sprintf("interval(%v,%v,%v)", Storex(i.from), Storex(i.to), Storex(i.by))
	}
	name := intervalStrategyName(i.strategy.id())
	return fmt.Sprintf("interval(%s,%s,%s,'%s')", Storex(i.from), Storex(i.to), Storex(i.by), name)
}

// Inpsect is part of Inspectable
func (i Interval) Inspect(n Inspection) {
	n.Properties["value"] = i.Value()
	n.Properties["length"] = (Int(i.to)-Int(i.from))/Int(i.by) + 1
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
	//reachedTo(i *Interval)
	//reachedFrom(i *Interval)
}

var intervalStrategies = map[string]intervalStrategy{
	"once":           strategyOnceFromTo{},
	"repeat":         strategyRepeatFromTo{},
	"two-way":        strategyOnceFromToFrom{},
	"repeat-two-way": strategyRepeatFromToFrom{},
}

type strategyOnceFromTo struct{}

func (s strategyOnceFromTo) id() int { return OnceFromTo }

type strategyRepeatFromTo struct{}

func (s strategyRepeatFromTo) id() int { return RepeatFromTo }

type strategyOnceFromToFrom struct{}

func (s strategyOnceFromToFrom) id() int { return OnceFromToFrom }

type strategyRepeatFromToFrom struct{}

func (s strategyRepeatFromToFrom) id() int { return RepeatFromToFrom }
