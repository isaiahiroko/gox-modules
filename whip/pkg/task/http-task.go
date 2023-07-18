package task

import (
	"github.com/gocraft/work"
)

type HttpTask struct {
	concurrency uint
	payload     TaskPayload
}

func (t *HttpTask) Name() string {
	return "HttpTask"
}

func (t *HttpTask) Args() map[string]any {
	return map[string]any{}
}

func (t *HttpTask) Options() work.JobOptions {
	return work.JobOptions{}
}

func (t *HttpTask) Action(w *work.WorkerPool) any {
	// logId := t.payload.Id + logSuffix
	// logger := logger.New(logId, os.Stdout)

	return func() {

		// if err != nil {
		// 	logger.Add(err.Error())
		// 	return
		// }
	}
}

func (t *HttpTask) Queue() string {
	return "Fast"
}

func (t *HttpTask) Concurrency() uint {
	return t.concurrency
}
