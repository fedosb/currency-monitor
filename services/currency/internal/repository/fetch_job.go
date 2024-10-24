package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/fedosb/currency-monitor/services/currency/internal/db/postgres"
	"github.com/fedosb/currency-monitor/services/currency/internal/entity"
)

type FetchJobRepository struct {
	db *postgres.DB
}

func NewFetchJobRepository(db *postgres.DB) *FetchJobRepository {
	return &FetchJobRepository{db: db}
}

type fetchJob struct {
	ID            int            `db:"id"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
	Name          string         `db:"name"`
	Status        string         `db:"status"`
	PlannedAt     time.Time      `db:"planned_at"`
	DownloadedAt  sql.NullTime   `db:"downloaded_at"`
	FailureReason sql.NullString `db:"failure_reason"`
	Retries       sql.NullInt64  `db:"retries"`
}

func (j fetchJob) Entity() entity.FetchJob {
	var downloadedAt *time.Time
	if j.DownloadedAt.Valid {
		downloadedAt = &j.DownloadedAt.Time
	}

	var failureReason *string
	if j.FailureReason.Valid {
		failureReason = &j.FailureReason.String
	}

	var retries *int
	if j.Retries.Valid {
		retries = new(int)
		*retries = int(j.Retries.Int64)
	}

	return entity.FetchJob{
		ID:            j.ID,
		CreatedAt:     j.CreatedAt,
		UpdatedAt:     j.UpdatedAt,
		Name:          j.Name,
		Status:        entity.ParseJobStatus(j.Status),
		PlannedAt:     j.PlannedAt,
		DownloadedAt:  downloadedAt,
		FailureReason: failureReason,
		Retries:       retries,
	}
}

func mapFetchJob(job entity.FetchJob) fetchJob {
	var downloadedAt sql.NullTime
	if job.DownloadedAt != nil {
		downloadedAt.Valid = true
		downloadedAt.Time = *job.DownloadedAt
	}

	var failureReason sql.NullString
	if job.FailureReason != nil {
		failureReason.Valid = true
		failureReason.String = *job.FailureReason
	}

	var retries sql.NullInt64
	if job.Retries != nil {
		retries.Valid = true
		retries.Int64 = int64(*job.Retries)
	}

	return fetchJob{
		ID:            job.ID,
		CreatedAt:     job.CreatedAt,
		UpdatedAt:     job.UpdatedAt,
		Name:          job.Name,
		Status:        job.Status.String(),
		PlannedAt:     job.PlannedAt,
		DownloadedAt:  downloadedAt,
		FailureReason: failureReason,
		Retries:       retries,
	}
}

type fetchJobList []fetchJob

func (l fetchJobList) Entities() []entity.FetchJob {
	jobs := make([]entity.FetchJob, 0, len(l))
	for job := range l {
		jobs = append(jobs, l[job].Entity())
	}
	return jobs
}

func (r *FetchJobRepository) Create(ctx context.Context, job entity.FetchJob) (entity.FetchJob, error) {
	jobModel := mapFetchJob(job)

	query, args, err := sqlx.Named(
		`INSERT INTO fetch_jobs (name, status, planned_at) 
				VALUES (:name, :status, :planned_at)
				RETURNING *`,
		&jobModel,
	)
	if err != nil {
		return entity.FetchJob{}, fmt.Errorf("mapping query: %w", err)
	}

	query = r.db.Rebind(query)

	if err = r.db.QueryRowxContext(ctx, query, args...).StructScan(&jobModel); err != nil {
		return entity.FetchJob{}, fmt.Errorf("query execution: %w", err)
	}

	return jobModel.Entity(), nil
}

func (r *FetchJobRepository) GetPlanned(ctx context.Context) ([]entity.FetchJob, error) {
	var jobs fetchJobList

	query := `UPDATE fetch_jobs SET 
            		status = 'PROCESSING',
            		updated_at = now()
            	WHERE status = 'READY' AND planned_at <= now() 
            	RETURNING *`

	if err := r.db.SelectContext(ctx, &jobs, query); err != nil {
		return nil, fmt.Errorf("query execution: %w", err)
	}

	return jobs.Entities(), nil
}

func (r *FetchJobRepository) GetFailed(ctx context.Context) ([]entity.FetchJob, error) {
	var jobs fetchJobList

	query := `UPDATE fetch_jobs SET 
            		status = 'PROCESSING',
            		updated_at = now()
            	WHERE status = 'FAILED' 
            	RETURNING *`

	if err := r.db.SelectContext(ctx, &jobs, query); err != nil {
		return nil, fmt.Errorf("query execution: %w", err)
	}

	return jobs.Entities(), nil
}

func (r *FetchJobRepository) Update(ctx context.Context, job entity.FetchJob) (entity.FetchJob, error) {
	jobModel := mapFetchJob(job)

	query, args, err := sqlx.Named(
		`UPDATE fetch_jobs SET
            	updated_at = now(),
            	name = :name,
				status = :status,
				planned_at = :planned_at,
				downloaded_at = :downloaded_at,
				failure_reason = :failure_reason,
				retries = :retries
			WHERE id = :id
			RETURNING *`,
		&jobModel,
	)
	if err != nil {
		return entity.FetchJob{}, fmt.Errorf("mapping query: %w", err)
	}

	query = r.db.Rebind(query)

	if err = r.db.QueryRowxContext(ctx, query, args...).StructScan(&jobModel); err != nil {
		return entity.FetchJob{}, fmt.Errorf("query execution: %w", err)
	}

	return jobModel.Entity(), nil
}
