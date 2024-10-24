-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rates (
    id          BIGSERIAL       PRIMARY KEY,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name        VARCHAR(255)    NOT NULL,
    date        DATE            NOT NULL,
    rate        NUMERIC(20, 15) NOT NULL,
    UNIQUE (name, date)
);

CREATE INDEX rates_name_date_idx ON rates(name, date);
CREATE INDEX rates_date_idx ON rates(date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS currency_rates;
-- +goose StatementEnd