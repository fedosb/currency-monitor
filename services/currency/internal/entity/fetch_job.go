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
