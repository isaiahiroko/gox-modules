# Whip

## Introduction

Whip is an HTTP-based scheduler. Its roughly based on

Once a request is successfully receive and persistent, it ensures the request gets to the destination eventually based on the provided config.

### Installation

1. Download the right binary for your platform from `assets` directory.
2. Run `install.sh` if your platform is linux/darwin

```bash
$ install.sh 0.0.1 darwin
```

### Usage & Commands

1. Run as a HTTP server

```
$ whip serve --port 25623
```

2. Run on CLI

```
$ makr run -f sample.json|sample.yaml|directory

The json or yaml file(s) must conform with the config struct below.
```

### Commands

1. `run`
1. `serve`

## Configuration

```go
package contract

type JobOptions struct {
	Timestamp          *int                `json:"timestamp,omitempty"` // Timestamp when the job was created. Default: Date.now()
	Priority           *int                `json:"priority,omitempty"`  // Ranges from 1 (highest priority) to MAX_INT (lowest priority).
	Delay              *int                `json:"delay,omitempty"`     // An amount of milliseconds to wait until this job can be processed. Default: 0
	Attempts           *int                `json:"attempts,omitempty"`  // The total number of attempts to try the job until it completes. Default: 0
	Backoff            *BackoffOptions     `json:"backoff,omitempty"`
	Lifo               *bool               `json:"lifo,omitempty"` // If true, adds the job to the right of the queue instead of the left (default false)
	RemoveOnComplete   *KeepJobsOptions    `json:"removeOnComplete,omitempty"`
	RemoveOnFail       *KeepJobsOptions    `json:"removeOnFail,omitempty"`
	KeepLogs           *int                `json:"keepLogs,omitempty"` // Maximum amount of log entries that will be preserved
	Repeat             *RepeatOptions      `json:"repeat,omitempty"`
	AutoRun            *bool               `json:"autorun,omitempty"`     // Condition to start processor at instance creation. Default: true
	Concurrency        *int                `json:"concurrency,omitempty"` // Amount of jobs that a single worker is allowed to work on in parallel. Default: 1
	Limiter            *RateLimiterOptions `json:"limiter,omitempty"`
	Metrics            *MetricsOptions     `json:"metrics,omitempty"`
	RemoveOnFailWorker *KeepJobsOptions    `json:"removeOnFailWorker,omitempty"`
}

type BackoffOptions struct {
	Type  string `json:"type,omitempty"`  // Name of the backoff strategy.
	Delay *int   `json:"delay,omitempty"` // Delay in milliseconds.
}

type KeepJobsOptions struct {
	Age   *int `json:"age,omitempty"`   // Maximum age in seconds for job to be kept.
	Count *int `json:"count,omitempty"` // Maximum count of jobs to be kept.
}

type RepeatOptions struct {
	Pattern     *string `json:"pattern,omitempty"`     // A repeat pattern
	Limit       *int    `json:"limit,omitempty"`       // Number of times the job should repeat at max
	Every       *int    `json:"every,omitempty"`       // Repeat after this amount of milliseconds
	Immediately *bool   `json:"immediately,omitempty"` // Repeated job should start right now
	Count       *int    `json:"count,omitempty"`       // The start value for the repeat iteration count
	PrevMillis  *int    `json:"prevMillis,omitempty"`
	Offset      *int    `json:"offset,omitempty"`
	JobID       *string `json:"jobId,omitempty"`
}

type RateLimiterOptions struct {
	Max      int `json:"max"`      // Max number of jobs to process in the time period specified in `duration`.
	Duration int `json:"duration"` // Time in milliseconds. During this time, a maximum of `max` jobs will be processed.
}

type MetricsOptions struct {
	MaxDataPoints *int `json:"maxDataPoints,omitempty"` // Enable gathering metrics for finished jobs. Output refers to all finished jobs, completed or failed.
}

type AdvancedOptions struct {
	// Define advanced options
}

```

## Request Status

- the status endpoint can also be called by clients to determine the status of any request.

### [License](./LICENSE.md)

### Todo
- complete http task action. action should take in worker/pool. subscriptions should be queued too.
- test http
- test cli
- add mongo store
- add mysql store