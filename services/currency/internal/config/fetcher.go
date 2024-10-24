package config

import "time"

type FetcherConfig interface {
	GetFetchInterval() time.Duration
	GetRequeueFailedJobsInterval() time.Duration
	GetJobRescheduleInterval() time.Duration
	GetFailedJobRescheduleInterval() time.Duration
	GetFailedJobMaxRetries() int
	GetSaveRatesMaxWorkers() int
}

type fetcherConfig struct {
	FetchInterval               time.Duration `env-required:"true" env:"FETCH_INTERVAL"`
	RequeueFailedJobsInterval   time.Duration `env-required:"true" env:"REQUEUE_FAILED_JOBS_INTERVAL"`
	JobRescheduleInterval       time.Duration `env-required:"true" env:"JOB_RESCHEDULE_INTERVAL"`
	FailedJobRescheduleInterval time.Duration `env-required:"true" env:"FAILED_JOB_RESCHEDULE_INTERVAL"`
	FailedJobMaxRetries         int           `env-required:"true" env:"FAILED_JOB_MAX_RETRIES"`
	SaveRatesMaxWorkers         int           `env-required:"true" env:"SAVE_RATES_MAX_WORKERS"`
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

func (c fetcherConfig) GetSaveRatesMaxWorkers() int {
	return c.SaveRatesMaxWorkers
}
