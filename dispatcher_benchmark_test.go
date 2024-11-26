package goevent

import (
	"context"
	"testing"
	"time"
)

func BenchmarkDispatcherHighPriority(b *testing.B) {
	dispatcher := NewDispatcher(20)
	defer dispatcher.Close()

	listener := &BenchmarkListener{}
	event := &TestEvent{
		name:     "high_event",
		priority: HighPriority,
	}

	dispatcher.Register(event, listener)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := dispatcher.Dispatch(context.Background(), event)
		if err != nil {
			b.Fatalf("failed to dispatch event: %v", err)
		}
	}
	dispatcher.Wait()
}

func BenchmarkDispatcherMediumPriority(b *testing.B) {
	dispatcher := NewDispatcher(20)
	defer dispatcher.Close()

	listener := &BenchmarkListener{}
	event := &TestEvent{
		name:     "med_event",
		priority: MedPriority,
	}

	dispatcher.Register(event, listener)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dispatcher.Dispatch(context.Background(), event)
	}
	dispatcher.Wait()
}

func BenchmarkDispatcherLowPriority(b *testing.B) {
	dispatcher := NewDispatcher(20)
	defer dispatcher.Close()

	listener := &BenchmarkListener{}
	event := &TestEvent{
		name:     "low_event",
		priority: LowPriority,
	}

	dispatcher.Register(event, listener)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dispatcher.Dispatch(context.Background(), event)
	}
	dispatcher.Wait()
}

type BenchmarkListener struct{}

func (l *BenchmarkListener) Handle(ctx context.Context, event Event) error {
	// Simulate some processing time
	time.Sleep(1 * time.Millisecond)
	return nil
}
