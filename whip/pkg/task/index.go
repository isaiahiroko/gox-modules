package task

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gocraft/work"

	"github.com/origine-run/whip/pkg/store"
)

type JSON map[string]any

type Link struct {
	Method  string `json:"method"`
	Url     string `json:"url"`
	Body    JSON   `json:"body"`
	Queries JSON   `json:"queries"`
	Headers JSON   `json:"headers"`
}

type BackoffOptions struct {
	Type  string `json:"type,omitempty"`  // Name of the backoff strategy.
	Delay int    `json:"delay,omitempty"` // Delay in milliseconds.
}

type KeepJobOptions struct {
	events []string
	Age    int `json:"age,omitempty"`   // Maximum age in seconds for job to be kept.
	Count  int `json:"count,omitempty"` // Maximum count of jobs to be kept.
}

type RepeatOptions struct {
	Pattern     string `json:"pattern,omitempty"`     // A repeat pattern
	Limit       int    `json:"limit,omitempty"`       // Number of times the job should repeat at max
	Every       int    `json:"every,omitempty"`       // Repeat after this amount of milliseconds
	Immediately bool   `json:"immediately,omitempty"` // Repeated job should start right now
}

type ConfigOptions struct {
	Timestamp   int64          `json:"timestamp"`
	Priority    int            `json:"priority"`
	Delay       int            `json:"delay"`
	Attempts    int            `json:"attempts"`
	Backoff     BackoffOptions `json:"backoff"`
	Remove      KeepJobOptions `json:"remove"`
	Repeat      RepeatOptions  `json:"repeat"`
	Concurrency int            `json:"concurrency"`
}

type SubscriptionOptions struct {
	OnStart    Link `json:"onStart"`
	OnRetry    Link `json:"onRetry"`
	OnProgress Link `json:"onProgress"`
	OnFail     Link `json:"onFail"`
	OnSuccess  Link `json:"onSuccess"`
	OnComplete Link `json:"onComplete"`
}

type TaskPayload struct {
	Id            string              `json:"id"`
	Checksum      string              `json:"checksum"`
	Force         bool                `json:"force"`
	Source        string              `json:"source"`
	Sink          Link                `json:"sink"`
	Subscriptions SubscriptionOptions `json:"subscriptions"`
	Config        ConfigOptions       `json:"config"`
}

type TaskContract interface {
	Queue() string
	Name() string
	Args() map[string]any
	Concurrency() uint
	Options() work.JobOptions
	Action(w *work.WorkerPool) any
}

var (
	defaultWorker = &worker{}
	Lists         = []TaskContract{
		&HttpTask{},
	}
	logSuffix = "log"
)

func Run() {
	// run
	tasks := []TaskContract{
		&HttpTask{},
	}
	defaultWorker.Run(tasks)
}

func Add(payload TaskPayload) ([]byte, error) {
	// idempotency
	checksumStore := store.New("checksum")
	exist := checksumStore.Exist(payload.Checksum)
	if exist {
		return nil, errors.New("The task already exist, use `Force` to run it if necessary")
	} else {
		checksumStore.Create(payload.Checksum, nil)
	}

	// store
	taskStore := store.New("task")
	err := taskStore.Create(payload.Id, payload)
	if err != nil {
		return nil, err
	}

	// add
	tasks := []TaskContract{
		&HttpTask{payload: payload, concurrency: uint(payload.Config.Concurrency)},
	}
	defaultWorker.Add(tasks)

	return json.Marshal(payload)
}

func GetOne(taskId string) ([]byte, error) {
	taskStore := store.New("task")

	task := taskStore.FindOne(taskId)

	logStore := store.New(taskId + "log")

	var logs []any

	err := logStore.FindMany(func(k string, v []byte) {
		var log any
		json.Unmarshal(v, &log)
		logs = append(logs, fmt.Sprintf("%s - %s", k, log.(string)))
	})

	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(struct {
		Task any
		Logs []any
	}{
		Task: task,
		Logs: logs,
	})

	return buf, err
}

func GetAll() ([]byte, error) {
	taskStore := store.New("task")
	var tasks []TaskPayload

	err := taskStore.FindMany(func(k string, v []byte) {
		var task TaskPayload
		json.Unmarshal(v, &task)
		tasks = append(tasks, task)
	})

	if err != nil {
		return nil, err
	}

	return json.Marshal(tasks)
}

func DeleteOne(taskId string) error {
	buf, err := GetOne(taskId)

	var v TaskPayload
	err = json.Unmarshal(buf, &v)
	if err != nil {
		return err
	}

	err = store.Empty(taskId + "log")
	if err != nil {
		return err
	}

	checksumStore := store.New("checksum")
	err = checksumStore.Delete(v.Checksum)
	if err != nil {
		return err
	}

	taskStore := store.New("task")
	return taskStore.Delete(taskId)
}

func DeleteAll() error {
	taskStore := store.New("task")

	err := taskStore.FindMany(func(k string, v []byte) {
		var task TaskPayload
		json.Unmarshal(v, &task)

		err := store.Empty(task.Id + "log")
		if err != nil {
			fmt.Println(err)
		}
	})

	if err != nil {
		return err
	}

	err = store.Empty("checksum")
	if err != nil {
		return err
	}

	return store.Empty("task")
}
