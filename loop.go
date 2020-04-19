package melrose

import (
	"fmt"
)

type Loop struct {
	Target    Sequenceable
	isRunning bool
	stopChan  chan bool
}

func (l Loop) Storex() string {
	if s, ok := l.Target.(Storable); ok {
		return fmt.Sprintf("loop(%s)", s.Storex())
	}
	return ""

}

func (l *Loop) IsRunning() bool {
	return l.isRunning
}

func (l *Loop) Start(d AudioDevice) *Loop {
	if l.isRunning || d == nil {
		return l
	}
	l.stopChan = make(chan bool)
	l.isRunning = true
	go func() {
		for {
			select {
			case <-l.stopChan:
				goto stop
			default:
				d.Play(l.Target)
			}
		}
	stop:
		l.isRunning = false
	}()
	return l
}

func (l *Loop) Stop() *Loop {
	if !l.isRunning {
		return l
	}
	l.stopChan <- true
	return l
}
