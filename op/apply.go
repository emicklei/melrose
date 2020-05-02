package op

import (
	"fmt"

	"github.com/emicklei/melrose"
)

type Apply struct {
	Target []interface{}
}

func (a Apply) Call(s melrose.Sequenceable) melrose.Sequenceable {
	result := s
	for _, each := range a.Target {
		result = a.withCall(each, result)
	}
	return result
}

func (a Apply) withCall(f interface{}, s melrose.Sequenceable) melrose.Sequenceable {
	if seq2seq, ok := f.(func(melrose.Sequenceable) melrose.Sequenceable); ok {
		fmt.Println("seq2seq")
		return seq2seq(s)
	}
	if intSlice2int, ok := f.(func(...interface{}) interface{}); ok {
		fmt.Println("intSlice2int")
		r := intSlice2int(s)
		if s, ok := r.(melrose.Sequenceable); ok {
			return s
		}
		fmt.Printf("r:%T\n", r)
		return nil
	}
	if int2int, ok := f.(func(interface{}) interface{}); ok {
		fmt.Println("int2int")
		r := int2int(s)
		if s, ok := r.(melrose.Sequenceable); ok {
			return s
		}
		fmt.Printf("r:%T\n", r)
		return nil
	}
	fmt.Printf("f:%T\n", f)
	return nil
}
