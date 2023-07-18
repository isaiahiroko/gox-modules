package dto

type TaskOptions struct {
	Timestamp        *int                `json:"timestamp,omitempty"` // Timestamp when the job was created. Default: Date.now()
	Priority         *int                `json:"priority,omitempty"`  // Ranges from 1 (highest priority) to MAX_INT (lowest priority).
	Delay            *int                `json:"delay,omitempty"`     // An amount of milliseconds to wait until this job can be processed. Default: 0
	Attempts         *int                `json:"attempts,omitempty"`  // The total number of attempts to try the job until it completes. Default: 0
	Backoff          *BackoffOptions     `json:"backoff,omitempty"`
	RemoveOnComplete *KeepJobsOptions    `json:"removeOnComplete,omitempty"`
	RemoveOnFail     *KeepJobsOptions    `json:"removeOnFail,omitempty"`
	Repeat           *RepeatOptions      `json:"repeat,omitempty"`
	Concurrency      *int                `json:"concurrency,omitempty"` // Amount of jobs that a single worker is allowed to work on in parallel. Default: 1
	Limiter          *RateLimiterOptions `json:"limiter,omitempty"`
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
