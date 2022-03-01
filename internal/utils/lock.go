package utils

import "sync"

type Lock interface {
	Unlock()
}

type RWUnlock struct {
	rw *sync.RWMutex
}

func (r *RWUnlock) Unlock() {
	r.rw.RUnlock()
}

type Unlock struct {
	rw *sync.RWMutex
}

func (r *Unlock) Unlock() {
	r.rw.Unlock()
}

type RWLock struct {
	lock sync.Mutex
	rw   map[string]*sync.RWMutex
}

func NewRWLock() *RWLock {
	return &RWLock{
		rw: map[string]*sync.RWMutex{},
	}
}

func (r *RWLock) Lock(key string) Lock {
	lock := r.getRWLock(key)
	lock.Lock()

	return &Unlock{rw: lock}
}

//func (r *RWLock) Unlock(key string) {
//	lock := r.getRWLock(key)
//	lock.Unlock()
//}

func (r *RWLock) RLock(key string) Lock {
	lock := r.getRWLock(key)
	lock.RLock()

	return &RWUnlock{rw: lock}
}

//func (r *RWLock) RUnlock(key string) {
//	lock := r.getRWLock(key)
//	lock.RUnlock()
//}

func (r *RWLock) getRWLock(key string) *sync.RWMutex {
	r.lock.Lock()
	r.lock.Unlock()

	_, ex := r.rw[key]
	if !ex {
		r.rw[key] = &sync.RWMutex{}
	}

	return r.rw[key]
}
