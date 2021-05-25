package test

import (
	"testing"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/stretchr/testify/assert"
)

const (
	interval = 100
)

type simpleWorker struct {
	Counter int
}

func (w *simpleWorker) DoWork() {
	w.Counter++
}

func TestCycleControllerCallsDoWorkEvery100ms(t *testing.T) {
	cycle := util.NewCycleController(interval)
	worker := &simpleWorker{}
	cycle.Subscribe(worker)
	cycle.Start()
	time.Sleep(1 * time.Second)
	cycle.Stop()
	assert.Equal(t, 10, worker.Counter) // was subscribed full time, has 10 for counter
}

func TestCycleControllerCallsMultipleDoWorkEvery100ms(t *testing.T) {
	cycle := util.NewCycleController(interval)
	worker := &simpleWorker{Counter: 0}
	workertwo := &simpleWorker{Counter: 0}

	cycle.Subscribe(worker)
	cycle.Subscribe(workertwo)
	cycle.Start()
	time.Sleep(1 * time.Second)
	cycle.Stop()
	assert.Equal(t, 10, worker.Counter)    // was subscribed full time, has 10 for counter
	assert.Equal(t, 10, workertwo.Counter) // was subscribed full time, has 10 for counter
}

func TestCycleControllerUnsubscribeStopsCallsWithoutResetingCycles(t *testing.T) {
	cycle := util.NewCycleController(interval)
	worker := &simpleWorker{Counter: 0}
	workertwo := &simpleWorker{Counter: 0}

	cycle.Subscribe(worker)
	cycle.Subscribe(workertwo)
	cycle.Start()
	time.Sleep(550 * time.Millisecond)
	cycle.Unsubscribe(workertwo)
	time.Sleep(500 * time.Millisecond)
	cycle.Stop()
	assert.Equal(t, 10, worker.Counter)   // was subscribed full time, has 10 for counter
	assert.Equal(t, 5, workertwo.Counter) // was subscribed half time, has 5 for counter
}

func TestCycleControllerSubscribeWhileExecutingWillContinue(t *testing.T) {
	cycle := util.NewCycleController(interval)
	worker := &simpleWorker{Counter: 0}
	workertwo := &simpleWorker{Counter: 0}

	cycle.Subscribe(worker)
	cycle.Start()
	// sleep until after the 5th tick but subscribe before 6th
	time.Sleep(510 * time.Millisecond)
	cycle.Subscribe(workertwo)
	// sleep until after 10th tick, but before 11th tick
	time.Sleep(510 * time.Millisecond)
	cycle.Stop()
	assert.Equal(t, 10, worker.Counter)   // was subscribed full time, has 10 for counter
	assert.Equal(t, 5, workertwo.Counter) // was subscribed half time, has 5 for counter
}

func TestCycleControllerSubscribeAndUnsubscribeWillNotIncrement(t *testing.T) {
	cycle := util.NewCycleController(interval)
	worker := &simpleWorker{Counter: 0}
	workertwo := &simpleWorker{Counter: 0}

	cycle.Subscribe(worker)
	cycle.Subscribe(workertwo)
	cycle.Unsubscribe(workertwo)
	cycle.Start()
	time.Sleep(1000 * time.Millisecond)
	cycle.Stop()
	assert.Equal(t, 10, worker.Counter)   // was subscribed full time, has 10 for counter
	assert.Equal(t, 0, workertwo.Counter) // was subscribed no time, has 0 for counter
}
