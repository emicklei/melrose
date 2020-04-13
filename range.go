package melrose

import "fmt"

type Interval struct {
	from  int
	to    int
	by    int
	value int
}

func (i *Interval) Value() interface{} {
	c := i.value
	next := c + i.by
	if c > i.to {
		i.value = i.from
	} else {
		i.value = next
	}
	return c
}

func NewInterval(from, to, by int) *Interval {
	return &Interval{from: from, to: to, by: by, value: from}
}

func (i Interval) Storex() string {
	return fmt.Sprintf("interval(%d,%d,%d)", i.from, i.to, i.by)
}
