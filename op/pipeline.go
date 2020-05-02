package op

import (
	"fmt"

	"github.com/emicklei/melrose"
)

type Pipeline struct {
	Target []interface{}
}

type pipelineStream struct {
	i      int
	target []interface{}
}

func (p *pipelineStream) atEnd() bool {
	return p.i == len(p.target)
}

func (p *pipelineStream) next() interface{} {
	n := p.target[p.i]
	p.i++
	return n
}

func (a Pipeline) Execute(s melrose.Sequenceable) (melrose.Sequenceable, error) {
	result := s
	stream := &pipelineStream{target: a.Target}
	for !stream.atEnd() {
		callResult, err := a.withExecute(stream, result)
		if err != nil {
			return nil, err
		}
		result = callResult
	}
	return result, nil
}

func (a Pipeline) withExecute(stream *pipelineStream, s melrose.Sequenceable) (melrose.Sequenceable, error) {
	f := stream.next()
	if seq2seq, ok := f.(func(melrose.Sequenceable) melrose.Sequenceable); ok {
		return seq2seq(s), nil
	}
	if intSlice2int, ok := f.(func(...interface{}) interface{}); ok {
		r := intSlice2int(s)
		s, ok := r.(melrose.Sequenceable)
		if !ok {
			return nil, fmt.Errorf("expected melrose.Sequenceable but got %v:(%T)", r, r)
		}
		return s, nil
	}
	if int2int, ok := f.(func(interface{}) interface{}); ok {
		r := int2int(s)
		s, ok := r.(melrose.Sequenceable)
		if !ok {
			return nil, fmt.Errorf("expected melrose.Sequenceable but got %v:(%T)", r, r)
		}
		return s, nil
	}
	if intint2int, ok := f.(func(int, interface{}) interface{}); ok {
		if stream.atEnd() {
			return nil, fmt.Errorf("expected int but out of arguments")
		}
		pop := stream.next()
		i, ok := pop.(int)
		if !ok {
			return nil, fmt.Errorf("expected int but got %v:(%T)", pop, pop)
		}
		r := intint2int(i, s)
		s, ok := r.(melrose.Sequenceable)
		if !ok {
			return nil, fmt.Errorf("expected melrose.Sequenceable but got %v:(%T)", r, r)
		}
		return s, nil
	}
	if stringint2int, ok := f.(func(string, interface{}) interface{}); ok {
		if stream.atEnd() {
			return nil, fmt.Errorf("expected string but out of arguments")
		}
		pop := stream.next()
		st, ok := pop.(string)
		if !ok {
			return nil, fmt.Errorf("expected string but got %v:(%T)", pop, pop)
		}
		r := stringint2int(st, s)
		s, ok := r.(melrose.Sequenceable)
		if !ok {
			return nil, fmt.Errorf("expected melrose.Sequenceable but got %v:(%T)", r, r)
		}
		return s, nil
	}
	return nil, fmt.Errorf("unhandled function|argument : %T", f)
}
