package melrose

import "fmt"

// Interval is a Valueable that has a Value between [from..to] and increments with [by].
// If the end of the interval is reached then the Value is set to [from].
// The fields of an Interval are also Valueable.
type Interval struct {
	from  Valueable
	to    Valueable
	by    Valueable
	value int
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

// NewInterval creates new Interval. The initial Value is set to [from].
func NewInterval(from, to, by Valueable) *Interval {
	start := Int(from)
	return &Interval{from: from, to: to, by: by, value: start}
}

// Storex is part of Storable.
func (i Interval) Storex() string {
	return fmt.Sprintf("interval(%v,%v,%v)", i.from, i.to, i.by)
}
