// This struct/interfaces combo sets up a process that executes on
// a periodic cycle.  Creating a CycleController allows you to subscribe/unsubscribe
// with a Worker interface with a DoWork method.  This DoWork method will be
// executed at the specified interval of the CycleController.

package util

import (
	"time"
)

type CycleController interface {
	Subscribe(Worker)
	Unsubscribe(Worker)
	Start()
	Stop()
}

type cycleController struct {
	interval     int
	workers      *[]*Worker
	timer        time.Ticker
	done         chan bool
	workerUpdate chan *[]*Worker
	processing   bool
}

type Worker interface {
	DoWork()
}

func NewCycleController(intervalMilliseconds int) CycleController {
	cycle := cycleController{
		interval:     intervalMilliseconds * int(time.Millisecond),
		workers:      &[]*Worker{},
		timer:        *time.NewTicker(time.Duration(intervalMilliseconds * int(time.Millisecond))),
		done:         make(chan bool),
		workerUpdate: make(chan *[]*Worker, 1),
		processing:   false,
	}

	return &cycle
}

func (c *cycleController) Subscribe(worker Worker) {
	workers := *c.workers
	workers = append(workers, &worker)

	if !c.processing {
		c.workers = &workers
	} else {
		c.workerUpdate <- &workers
	}
}

func (c *cycleController) Unsubscribe(worker Worker) {
	workers := *c.workers
	index := len(workers)

	for i, existing := range workers {
		if *existing == worker {
			index = i

			break
		}
	}

	copy(workers[index:], workers[index+1:])
	workers[len(workers)-1] = nil
	workers = workers[:len(workers)-1]

	if !c.processing {
		c.workers = &workers
	} else {
		c.workerUpdate <- &workers
	}
}

func (c *cycleController) Start() {
	c.timer.Reset(time.Duration(c.interval))
	c.processing = true

	go func(t *time.Ticker, d *chan bool, workers *[]*Worker, workersUpdate *chan *[]*Worker) {
		currWorkers := workers

		for {
			select {
			case nw := <-*workersUpdate:
				currWorkers = nw
			case <-*d:
				return
			case <-t.C:
				for _, worker := range *currWorkers {
					if worker != nil {
						w := *worker
						w.DoWork()
					}
				}
			}
		}
	}(&c.timer, &c.done, c.workers, &c.workerUpdate)
}

func (c *cycleController) Stop() {
	c.done <- true
	c.timer.Stop()
	c.processing = false
}
