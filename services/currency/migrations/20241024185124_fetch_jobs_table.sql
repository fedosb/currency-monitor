-- +goose Up
-- +goose StatementBegin
CREATE TYPE fetch_job_status AS ENUM ('READY', 'PROCESSING', 'FAILED');

CREATE TABLE fetch_jobs (
    id                  BIGSERIAL           PRIMARY KEY,
    created_at          TIMESTAMP           NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP           NOT NULL DEFAULT NOW(),
    name                VARCHAR(32)         NOT NULL UNIQUE,
    status              fetch_job_status    NOT NULL,
    planned_at          TIMESTAMP           NOT NULL,
    downloaded_at       TIMESTAMP,
    failure_reason      TEXT,
    retries             INT
);

CREATE INDEX fetch_jobs_planned_at_idx ON fetch_jobs(planned_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fetch_jobs;

DROP TYPE fetch_job_status;
-- +goose StatementEnd
