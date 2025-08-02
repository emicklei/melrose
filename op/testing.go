package op

import (
	"errors"

	"github.com/emicklei/melrose/core"
)

type failingNoteConvertable struct{}

func (f failingNoteConvertable) ToNote() (core.Note, error) {
	return core.Rest4, errors.New("i am a failure")
}

func (f failingNoteConvertable) S() core.Sequence {
	return core.EmptySequence
}

func (f failingNoteConvertable) Storex() string {
	return ""
}
