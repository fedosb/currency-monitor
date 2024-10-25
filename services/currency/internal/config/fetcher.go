package config

import (
	"strings"
	"time"
)

type FetcherConfig interface {
	GetFetchInterval() time.Duration
	GetRequeueFailedJobsInterval() time.Duration
	GetJobRescheduleInterval() time.Duration
	GetFailedJobRescheduleInterval() time.Duration
	GetFailedJobMaxRetries() int
	GetProcessJobMaxWorkers() int
	GetCurrencies() []string
	GetRunImmediately() bool
}

type fetcherConfig struct {
	FetchInterval               time.Duration `env-required:"true" env:"FETCH_INTERVAL"`
	RequeueFailedJobsInterval   time.Duration `env-required:"true" env:"REQUEUE_FAILED_JOBS_INTERVAL"`
	JobRescheduleInterval       time.Duration `env-required:"true" env:"JOB_RESCHEDULE_INTERVAL"`
	FailedJobRescheduleInterval time.Duration `env-required:"true" env:"FAILED_JOB_RESCHEDULE_INTERVAL"`
	FailedJobMaxRetries         int           `env-required:"true" env:"FAILED_JOB_MAX_RETRIES"`
	ProcessJobMaxWorkers        int           `env-required:"true" env:"PROCESS_JOB_MAX_WORKERS"`
	Currencies                  string        `env:"CURRENCIES"`
	RunImmediately              bool          `env:"RUN_IMMEDIATELY"`
}

func (c fetcherConfig) GetFetchInterval() time.Duration {
	return c.FetchInterval
}

func (c fetcherConfig) GetRequeueFailedJobsInterval() time.Duration {
	return c.RequeueFailedJobsInterval
}

func (c fetcherConfig) GetJobRescheduleInterval() time.Duration {
	return c.JobRescheduleInterval
}

func (c fetcherConfig) GetFailedJobRescheduleInterval() time.Duration {
	return c.FailedJobRescheduleInterval
}

func (c fetcherConfig) GetFailedJobMaxRetries() int {
	return c.FailedJobMaxRetries
}

func (c fetcherConfig) GetProcessJobMaxWorkers() int {
	return c.ProcessJobMaxWorkers
}

func (c fetcherConfig) GetCurrencies() []string {
	return strings.Split(c.Currencies, ",")
}

func (c fetcherConfig) GetRunImmediately() bool {
	return c.RunImmediately
}
