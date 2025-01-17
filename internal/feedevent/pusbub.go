package feedevent

import (
	"sync"
)

type PubSub struct {
	subs       map[chan interface{}]bool
	bufferSize int
	closed     bool

	m sync.RWMutex
}

func NewPubSub(bufferSize int) *PubSub {
	return &PubSub{
		subs:       make(map[chan interface{}]bool),
		bufferSize: bufferSize,
	}
}

func (p *PubSub) Subscribe() chan interface{} {
	p.m.Lock()
	defer p.m.Unlock()

	ch := make(chan interface{}, p.bufferSize)
	p.subs[ch] = true
	return ch
}

func (p *PubSub) Publish(msg interface{}) {
	p.m.RLock()
	defer p.m.RUnlock()

	if p.closed {
		return
	}
	for ch := range p.subs {
		ch <- msg
	}
}

func (p *PubSub) PublishNoWait(msg interface{}) {
	p.m.RLock()
	defer p.m.RUnlock()

	if p.closed {
		return
	}

	for ch := range p.subs {
		select {
		case ch <- msg:
		default:
		}
	}
}

func (p *PubSub) Unsubscribe(ch chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	if _, ok := p.subs[ch]; !ok {
		return
	}

	delete(p.subs, ch)
}

func (p *PubSub) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if !p.closed {
		p.closed = true
		for ch := range p.subs {
			delete(p.subs, ch)
			close(ch)
		}
	}
}
