package models

import "sync"

type Whitelist struct {
	data map[string]struct{}
	mx   sync.RWMutex
}

func (w *Whitelist) Put(values ...string) {
	w.mx.Lock()
	defer w.mx.Unlock()

	if w.data == nil {
		w.data = make(map[string]struct{})
	}

	for _, v := range values {
		w.data[v] = struct{}{}
	}
}

func (w *Whitelist) Has(value string) bool {
	w.mx.RLock()
	_, ok := w.data[value]
	w.mx.RUnlock()

	return ok
}
