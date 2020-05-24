package op
import "fmt"
import . "github.com/emicklei/melrose"

type Parallel struct {
	Target Sequenceable
}

func (p Parallel) S() Sequence {
	n := []Note{}
	p.Target.S().NotesDo(func(each Note) {
		n = append(n, each)
	})
	return Sequence{Notes: [][]Note{n}}
}

func (p Parallel) Storex() string {
	if s, ok := p.Target.(Storable); ok {
		return fmt.Sprintf("parallel(%s)", s.Storex())
	}
	return "" 
}