package task

import (
	"os"
	"os/signal"

	"github.com/gocraft/work"
)

type worker struct{}

func (w *worker) Add(tasks []TaskContract) {
	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		queue := &queue{name: task.Queue()}
		queue.Add(task.Name(), task.Args())
	}
}

func (w *worker) Run(tasks []TaskContract) {
	var pool *work.WorkerPool
	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		pool := work.NewWorkerPool(task, task.Concurrency(), task.Queue(), redisPool)
		pool.JobWithOptions(task.Name(), task.Options(), task.Action(pool))
	}

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}
