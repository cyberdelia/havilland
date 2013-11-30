// Generic single flight
package havilland

import (
	"sync"
)

type Airline struct {
	flights map[interface{}]*flight
	m sync.Mutex
}

type flight struct {
	value interface{}
	err   error
	sync.WaitGroup
}

func (a *Airline) Fly(identifier interface{}, fn func() (interface{}, error)) (interface{}, error) {
	// Lock airspace
	a.m.Lock()

	// Initialize flights
	if a.flights == nil {
		a.flights = make(map[interface{}]*flight)
	}

	// If an existing flight exists reuse it
	if f, ok := a.flights[identifier]; ok {
		a.m.Unlock()
		f.Wait()
		return f.value, f.err
	}

	// Create a new flight
	f := new(flight)
	f.Add(1)
	a.flights[identifier] = f
	a.m.Unlock()

	// Execute
	f.value, f.err = fn()
	f.Done()

	// Remove flight from airspace
	a.m.Lock()
	delete(a.flights, identifier)
	a.m.Unlock()

	return f.value, f.err
}
