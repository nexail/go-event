package goevent

import (
	"context"
	"log"
	"sync"
)

type Dispatcher struct {
	listeners sync.Map
	wg        sync.WaitGroup
	highJobs  chan Event
	medJobs   chan Event
	lowJobs   chan Event
	workers   int
	shutdown  chan struct{} // Channel to signal shutdown
}

func NewDispatcher(workers int) *Dispatcher {
	d := &Dispatcher{
		highJobs: make(chan Event, 1000), // Buffered channel for high-priority events
		medJobs:  make(chan Event, 1000), // Buffered channel for medium-priority events
		lowJobs:  make(chan Event, 1000), // Buffered channel for low-priority events
		workers:  workers,                // Number of workers
		shutdown: make(chan struct{}),    // Channel to signal shutdown
	}
	d.startWorkers()
	return d
}

func (d *Dispatcher) startWorkers() {
	for i := 0; i < d.workers; i++ {
		go d.worker()
	}
}

func (d *Dispatcher) worker() {
	for {
		select {
		case event, ok := <-d.highJobs:
			if !ok {
				return // Channel closed, worker should stop
			}
			d.handleEvent(event)
		case event, ok := <-d.medJobs:
			if !ok {
				return // Channel closed, worker should stop
			}
			d.handleEvent(event)
		case event, ok := <-d.lowJobs:
			if !ok {
				return // Channel closed, worker should stop
			}
			d.handleEvent(event)
		case <-d.shutdown:
			return // Stop the worker
		}
	}
}

func (d *Dispatcher) handleEvent(event Event) {
	if err := d.processEvent(context.Background(), event); err != nil {
		log.Printf("Error processing event %s: %v", event.Name(), err)
	}
}

func (d *Dispatcher) processEvent(ctx context.Context, event Event) error {
	if event == nil {
		return nil // Exit if event is nil to prevent dereferencing
	}

	if listeners, found := d.listeners.Load(event.Name()); found {
		for _, listener := range listeners.([]Listener) {
			if listener == nil {
				continue // Skip if listener is nil
			}
			d.wg.Add(1)
			go func(l Listener) {
				defer d.wg.Done()
				if err := l.Handle(ctx, event); err != nil {
					log.Printf("Error handling event %s: %v", event.Name(), err)
				}
			}(listener)
		}
	}
	return nil
}

func (d *Dispatcher) Register(event Event, listeners ...Listener) {
	for _, listener := range listeners {
		existingListeners, _ := d.listeners.LoadOrStore(event.Name(), []Listener{})
		d.listeners.Store(event.Name(), append(existingListeners.([]Listener), listener))
	}
}

func (d *Dispatcher) Unregister(event Event, listener Listener) {
	if listeners, found := d.listeners.Load(event.Name()); found {
		updatedListeners := listeners.([]Listener)
		for i, l := range updatedListeners {
			if l == listener {
				d.listeners.Store(event.Name(), append(updatedListeners[:i], updatedListeners[i+1:]...))
				break
			}
		}
	}
}

// count the number of listeners
func (d *Dispatcher) Listeners() int {
	count := 0
	d.listeners.Range(func(_, value interface{}) bool {
		count += len(value.([]Listener))
		return true
	})
	return count
}

func (d *Dispatcher) Dispatch(ctx context.Context, event Event) error {
	switch event.Priority() {
	case HighPriority:
		d.highJobs <- event
	case MedPriority:
		d.medJobs <- event
	default:
		d.lowJobs <- event
	}
	return nil
}

func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

// cleanup the dispatcher
func (d *Dispatcher) Close() {
	close(d.shutdown) // Signal all workers to stop
	d.Wait()          // Wait for all ongoing events to be processed
	close(d.highJobs)
	close(d.medJobs)
	close(d.lowJobs)
}
