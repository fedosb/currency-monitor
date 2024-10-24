package entity

import "time"

type FetchJob struct {
	ID            int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Status        JobStatus
	PlannedAt     time.Time
	DownloadedAt  *time.Time
	FailureReason *string
	Retries       *int
}

func (j *FetchJob) Succeed() {
	now := time.Now().UTC()

	j.Status = JobStatusReady
	j.DownloadedAt = &now
	j.FailureReason = nil
	j.Retries = nil
}

func (j *FetchJob) Reschedule(interval time.Duration) {
	j.PlannedAt = time.Now().UTC().Add(interval).Truncate(time.Second)
}

func (j *FetchJob) Fail(reason string) {
	j.Status = JobStatusFailed
	j.FailureReason = &reason
	if j.Retries == nil {
		j.Retries = new(int)
	}
}

type JobStatus string

const (
	JobStatusUnknown    JobStatus = "UNKNOWN"
	JobStatusReady      JobStatus = "READY"
	JobStatusProcessing JobStatus = "PROCESSING"
	JobStatusFailed     JobStatus = "FAILED"
)

func (s JobStatus) String() string {
	return string(s)
}

func ParseJobStatus(s string) JobStatus {
	switch s {
	case JobStatusReady.String():
		return JobStatusReady
	case JobStatusProcessing.String():
		return JobStatusProcessing
	case JobStatusFailed.String():
		return JobStatusFailed
	default:
		return JobStatusUnknown
	}
}
