package world

import (
	"sync"
)

type ID uint64

type World struct {
	nextID  ID
	idMutex sync.Mutex
}

func New() *World {
	return &World{}
}

func (w *World) NextID() ID {
	w.idMutex.Lock()
	defer w.idMutex.Unlock()
	id := w.nextID
	w.nextID += 1
	return id
}
