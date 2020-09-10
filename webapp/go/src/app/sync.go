package main

import "sync"

type AutoIncrement struct {
	Current int
	Mux     sync.Mutex
}

var autoIncrement = &AutoIncrement{Current: 0}

func (ai *AutoIncrement) FetchID() int {
	ai.Mux.Lock()
	defer ai.Mux.Unlock()
	ai.Current += 1

	return ai.Current
}
