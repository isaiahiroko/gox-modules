package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/origine-run/makr/pkg/container"
	"github.com/origine-run/makr/pkg/logger"
	"github.com/origine-run/makr/pkg/repository"
	"github.com/origine-run/makr/pkg/store"
)

type JobModel struct {
	Id             string `json:"id"`
	GitHost        string `json:"git-host"`
	GitUsername    string `json:"git-username"`
	GitRepo        string `json:"git-repo"`
	GitPassword    string `json:"git-password"`
	GitRemote      string `json:"git-remote"`
	GitBranch      string `json:"git-branch"`
	DockerHost     string `json:"docker-host"`
	DockerUsername string `json:"docker-username"`
	DockerRegistry string `json:"docker-registry"`
	DockerPassword string `json:"docker-password"`
	ImageVersion   string `json:"image-version"`
	Force          bool   `json:"force"`
	Checksum       string `json:"checksum"`
}

var (
	logSuffix = "log"
)

type Job struct {
	j JobModel
}

func (j *Job) Run(payload JobModel) {
	logId := payload.Id + logSuffix
	logger := logger.New(logId, os.Stdout)

	// repo
	repository := repository.New(fmt.Sprintf("%s/%s/%s", payload.GitHost, payload.GitUsername, payload.GitRepo), payload.GitPassword, payload.GitRemote, payload.GitBranch, logger)
	err := repository.Clone(fmt.Sprintf("./tmp/%s/%s", payload.GitUsername, payload.GitRepo))
	if err != nil {
		logger.Add(err.Error())
		return
	}

	// r.Checkout(gCommit)

	// container
	c := container.New(logger)
	err = c.Build(
		fmt.Sprintf("./tmp/%s/%s", payload.GitUsername, payload.GitRepo),
		fmt.Sprintf("%s/%s:%s", payload.DockerUsername, payload.DockerRegistry, payload.ImageVersion),
	)
	if err != nil {
		logger.Add(err.Error())
		return
	}

	err = c.Push(payload.DockerUsername, payload.DockerPassword, payload.DockerHost, fmt.Sprintf("%s/%s:%s", payload.DockerUsername, payload.DockerRegistry, payload.ImageVersion))
	if err != nil {
		logger.Add(err.Error())
		return
	}
}

func (j *Job) Add(payload JobModel) ([]byte, error) {
	// idempotency
	checksumStore := store.New("checksum")
	exist := checksumStore.Exist(payload.Checksum)
	if exist {
		return nil, errors.New("The job already exist, use `Force` to run it if necessary")
	} else {
		checksumStore.Create(payload.Checksum, nil)
	}

	// store
	jobStore := store.New("job")
	err := jobStore.Create(payload.Id, payload)
	if err != nil {
		return nil, err
	}

	// run
	go j.Run(payload)

	return json.Marshal(payload)
}

func (j *Job) GetOne(jobId string) ([]byte, error) {
	jobStore := store.New("job")

	job := jobStore.FindOne(jobId)

	logStore := store.New(jobId + "log")

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
		Job  any
		Logs []any
	}{
		Job:  job,
		Logs: logs,
	})

	return buf, err
}

func (j *Job) GetAll() ([]byte, error) {
	jobStore := store.New("job")
	var jobs []JobModel

	err := jobStore.FindMany(func(k string, v []byte) {
		var job JobModel
		json.Unmarshal(v, &job)
		jobs = append(jobs, job)
	})

	if err != nil {
		return nil, err
	}

	return json.Marshal(jobs)
}

func (j *Job) DeleteOne(jobId string) error {
	buf, err := j.GetOne(jobId)

	var v JobModel
	err = json.Unmarshal(buf, &v)
	if err != nil {
		return err
	}

	err = store.Empty(jobId + "log")
	if err != nil {
		return err
	}

	checksumStore := store.New("checksum")
	err = checksumStore.Delete(v.Checksum)
	if err != nil {
		return err
	}

	jobStore := store.New("job")
	return jobStore.Delete(jobId)
}

func (j *Job) DeleteAll() error {
	jobStore := store.New("job")

	err := jobStore.FindMany(func(k string, v []byte) {
		var job JobModel
		json.Unmarshal(v, &job)

		err := store.Empty(job.Id + "log")
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

	return store.Empty("job")
}

func New() *Job {
	return &Job{}
}
