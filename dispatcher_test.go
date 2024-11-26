package goevent

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DispatcherTestSuite struct {
	suite.Suite
	dispatcher   *Dispatcher
	mockListener *MockListener
}

type MockListener struct {
	mock.Mock
}

func (m *MockListener) Handle(ctx context.Context, event Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func TestDispatcherTestSuite(t *testing.T) {
	suite.Run(t, new(DispatcherTestSuite))
}

func (suite *DispatcherTestSuite) SetupTest() {
	suite.dispatcher = NewDispatcher(20)
	suite.mockListener = new(MockListener)
}

func (suite *DispatcherTestSuite) TearDownTest() {
	suite.dispatcher.Close()
}

func (suite *DispatcherTestSuite) TestRegisterAndDispatchHighPriorityEvent() {
	var wg sync.WaitGroup
	wg.Add(1)

	event := &TestEvent{
		name:     "high_event",
		priority: HighPriority,
	}

	suite.mockListener.On("Handle", mock.Anything, event).Run(func(args mock.Arguments) {
		defer wg.Done() // Mark event as handled
	}).Return(nil)

	suite.dispatcher.Register(event, suite.mockListener)
	err := suite.dispatcher.Dispatch(context.Background(), event)

	suite.NoError(err)

	// Wait for the event to be handled
	wg.Wait()

	// Assert that the listener was called
	suite.mockListener.AssertCalled(suite.T(), "Handle", mock.Anything, event)
}

func (suite *DispatcherTestSuite) TestRegisterAndDispatchMediumPriorityEvent() {
	var wg sync.WaitGroup
	wg.Add(1)

	event := &TestEvent{
		name:     "med_event",
		priority: MedPriority,
	}

	suite.mockListener.On("Handle", mock.Anything, event).Run(func(args mock.Arguments) {
		defer wg.Done() // Mark event as handled
	}).Return(nil)

	suite.dispatcher.Register(event, suite.mockListener)
	err := suite.dispatcher.Dispatch(context.Background(), event)

	suite.NoError(err)

	// Wait for the event to be handled
	wg.Wait()

	// Assert that the listener was called
	suite.mockListener.AssertCalled(suite.T(), "Handle", mock.Anything, event)
}

func (suite *DispatcherTestSuite) TestRegisterAndDispatchLowPriorityEvent() {
	var wg sync.WaitGroup
	wg.Add(1)

	event := &TestEvent{
		name:     "low_event",
		priority: LowPriority,
	}

	suite.mockListener.On("Handle", mock.Anything, event).Run(func(args mock.Arguments) {
		defer wg.Done() // Mark event as handled
	}).Return(nil)

	suite.dispatcher.Register(event, suite.mockListener)
	err := suite.dispatcher.Dispatch(context.Background(), event)

	suite.NoError(err)

	// Wait for the event to be handled
	wg.Wait()

	// Assert that the listener was called
	suite.mockListener.AssertCalled(suite.T(), "Handle", mock.Anything, event)
}

type TestEvent struct {
	name     EventName
	priority EventPriority
}

func (e *TestEvent) Name() EventName {
	return e.name
}

func (e *TestEvent) Priority() EventPriority {
	return e.priority
}
