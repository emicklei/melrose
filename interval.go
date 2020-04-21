package melrose

import "fmt"

const (
	OnceFromTo = iota
	OnceFromToFrom
	RepeatFromTo
	RepeatFromToFrom
)

// Interval is a Valueable that has a Value between [from..to] and increments with [by].
// If the end of the interval is reached then the Value is set to [from].
// The fields of an Interval are also Valueable.
type Interval struct {
	from     Valueable
	to       Valueable
	by       Valueable
	strategy intervalStrategy
	value    int
}

// Value returns the current value of the interval and increases its value with [by].
func (i *Interval) Value() interface{} {
	c := i.value
	by := Int(i.by)
	next := c + by
	if by < 0 {
		if next < Int(i.from) {
			i.value = Int(i.to)
			return c
		}
	}
	if by > 0 {
		if next > Int(i.to) {
			i.value = Int(i.from)
			return c
		}
	}
	i.value = next
	return c
}

// NewInterval creates new Interval. The initial Value is set to [from]. Specify the repeat strategy.
func NewInterval(from, to, by Valueable, strategy int) *Interval {
	start := Int(from)
	return &Interval{from: from, to: to, by: by, value: start, strategy: asIntervalStrategy(strategy)}
}

// Storex is part of Storable.
func (i Interval) Storex() string {
	if i.strategy.id() == OnceFromTo {
		return fmt.Sprintf("interval(%v,%v,%v)", i.from, i.to, i.by)
	}
	name := IntervalStrategy(i.strategy.id())
	return fmt.Sprintf("interval(%v,%v,%v,'%s')", i.from, i.to, i.by, name)
}

func ParseIntervalStrategy(s string) intervalStrategy {
	if is, ok := intervalStrategies[s]; ok {
		return is
	}
	return strategyOnceFromTo{}
}
func IntervalStrategy(i int) string {
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
