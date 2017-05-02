package neko

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// Simple tuple around a test description and the work
type testModern struct {
	Name string
	Func func(t *testing.T)
}

// Keeps track of mocks, setups, and tests so they can be
// later coordinated.

type OrganizerModern struct {
	t *testing.T

	only    *testModern
	mocks   []*mock.Mock
	setup   []func()
	cleanup []func()
	tests   []testModern
}

// Create a new OrganizerModern against testing's T interface
func Modern(t *testing.T) *OrganizerModern {
	return &OrganizerModern{t: t}
}

// Track a github.com/stretchr/testify/mock.Mock along with the tests
func (o *OrganizerModern) CheckMock(m *mock.Mock) {
	o.mocks = append(o.mocks, m)
}

// Add some work to be done before each test
func (o *OrganizerModern) Setup(f func()) {
	o.setup = append(o.setup, f)
}

// Add some work to be done after each test
func (o *OrganizerModern) Cleanup(f func()) {
	o.cleanup = append(o.cleanup, f)
}

// Add a test.
func (o *OrganizerModern) It(name string, f func(t *testing.T)) {
	o.tests = append(o.tests, testModern{name, f})
}

func (o *OrganizerModern) Only(name string, f func(t *testing.T)) {
	o.only = &testModern{name, f}
}

// Useful by allowing the developer to simply add 'N' before
// It to disable a block.
func (o *OrganizerModern) NIt(name string, f func(t *testing.T)) {
	o.tests = append(o.tests, testModern{name, nil})
}

// Coordinate running the tests with the setups and mocks
func (o *OrganizerModern) Run() {
	if o.only != nil {
		o.runTest(o.only)
		return
	}

	for _, test := range o.tests {
		o.runTest(&test)
	}
}

func (o *OrganizerModern) runTest(test *testModern) {
	if test.Func == nil {
		o.t.Logf("==== DISABLED: %s ====", test.Name)
		return
	}

	o.t.Run(test.Name, func(t *testing.T) {
		for _, mock := range o.mocks {
			mock.ExpectedCalls = nil
			mock.Calls = nil
		}

		for _, setup := range o.setup {
			setup()
		}

		defer o.runCleanup()

		test.Func(t)

		for _, mock := range o.mocks {
			mock.AssertExpectations(o.t)
		}
	})
}

func (o *OrganizerModern) runCleanup() {
	for _, cleanup := range o.cleanup {
		cleanup()
	}
}

// Have fun with neko!
func (o *OrganizerModern) Meow() {
	o.Run()
}
