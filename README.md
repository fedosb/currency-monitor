# Currency Monitor
**Project Description for Monitoring Currency Exchange Rates**

The project monitors currency exchange rates relative to the Russian ruble and consists of two services: **currency** and **gateway**.
It automates the collection of current currency rates with minimal latency and error handling to ensure data consistency even if there are connectivity issues with external APIs.

### Currency Service
The **currency** service manages currency rate monitoring with jobs scheduled for each currency. Each job has a status and a `planned_at` parameter indicating the scheduled execution time.

**Process Workflow:**
- A cron job triggers all scheduled jobs, checking if their execution time matches the current moment.
- For jobs scheduled to run, a request is sent to the currency API to fetch the latest exchange rates.
- The retrieved data is saved to the database, the job is marked as complete, and it is rescheduled for a configurable interval.
- If a job fails to complete, its status is set to `FAILED`. A separate cron monitors such jobs, changes their status to `READY`, and increments the retry count if it hasn’t reached the maximum. Otherwise, these jobs are rescheduled according to the main interval.

### Gateway Service
The **gateway** service handles authentication and interaction with the **currency** service, supporting the following requests:
- User sign-in.
- Requests to fetch currency rates by name and date.
- Requests to fetch currency rates by name and a date range.

**How to Run the Project**

1. Navigate to the `deployment/local` directory and start the infrastructure using the `docker-compose-infrastructure.yml` file (this starts the database).
    - Command:
      ```bash
      docker-compose -f docker-compose-infrastructure.yml up -d
      ```

2. Then, start the services using the `docker-compose.yml` file. Docker images for the services are available on Docker Hub (authorization service) and GHCR (currency and gateway services).
    - Command:
      ```bash
      docker-compose -f docker-compose.yml up -d
      ```
**Currency Service Configuration**

The **currency** service has several environment variables to configure network, API, database, and job management settings:

- **Network Configuration**
    - `GRPC_ADDRESS`: The gRPC server address (default `:50051`).

- **Currency API Configuration**
    - `CURRENCY_API_URL`: The URL to fetch currency exchange rates (`https://latest.currency-api.pages.dev/v1/currencies`).

- **Database Configuration**
    - Basic settings for connecting to the database, such as host, port, username, password, and database name.
    - Additional parameters control SSL (`DB_SSL`), connection limits (`DB_MAX_OPEN_CONNECTIONS`, `DB_MAX_IDLE_CONNECTIONS`, `DB_MAX_CONNECTION_LIFETIME`), and whether to apply migrations on start-up (`DB_APPLY_MIGRATIONS`).

- **Job Management Configuration**
    - `FETCH_INTERVAL`: Interval for fetching currency data (e.g., `15s`).
    - `REQUEUE_FAILED_JOBS_INTERVAL`: Interval for checking failed jobs to requeue (`30s`).
    - `JOB_RESCHEDULE_INTERVAL`: Default interval for rescheduling completed jobs (`24h`).
    - `FAILED_JOB_RESCHEDULE_INTERVAL`: Interval to retry failed jobs (`5m`).
    - `FAILED_JOB_MAX_RETRIES`: Maximum retry count for failed jobs (`3`).
    - `PROCESS_JOB_MAX_WORKERS`: Maximum workers to process jobs concurrently (`5`).
    - `CURRENCIES`: List of currencies to monitor (e.g., `usd,eur,gbp,algo`).
    - `RUN_IMMEDIATELY`: Whether to start the cron for monitoring jobs right away (`true`).

### Endpoints

#### Sign In
Endpoint for user authentication.

**Request:**
```bash
curl --location 'http://localhost:8081/api/sign-in' \
--header 'Content-Type: application/json' \
--data '{
    "login": "fedor",
    "password": "qwerty"
}'
```

**Responses:**
- **200 OK**
  ```json
  {
      "data": {
          "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6ImZlZG9yIiwiZXhwIjoxNzMwMTEzNDAxfQ.ReSuNOfdW4byfSlOrg7yy6w1MLZB6BMqPzMwhLe0V14"
      }
  }
  ```
- **401 Unauthorized**
  ```json
  {
      "error": "invalid login or password"
  }
  ```
- **500 Internal Server Error**
  ```json
  {
      "error": "something went wrong"
  }
  ```

---

#### Get Currency Rate by Name and Date
Retrieve the exchange rate for a specified currency on a particular date.

**Request:**
```bash
curl --location 'http://localhost:8081/api/currency/usd?date=2024-10-27' \
--header 'Authorization: Bearer Token'
```

**Responses:**
- **200 OK**
  ```json
  {
      "data": {
          "rate": {
              "id": 1,
              "created_at": "2024-10-27T21:34:11.186885Z",
              "updated_at": "2024-10-27T21:34:11.186885Z",
              "name": "usd",
              "date": "2024-10-27T00:00:00Z",
              "rate": 0.010286624
          }
      }
  }
  ```
- **400 Bad Request**
  ```json
  {
      "error": "validation error: correct date is required"
  }
  ```
- **401 Unauthorized**
  ```json
  {
      "error": "token expired or invalid"
  }
  ```
- **404 Not Found**
  ```json
  {
      "error": "not found"
  }
  ```
- **500 Internal Server Error**
  ```json
  {
      "error": "something went wrong"
  }
  ```

---

#### Get Currency Rates by Name and Date Range
Retrieve the exchange rates for a specified currency within a given date range.

**Request:**
```bash
curl --location 'http://localhost:8081/api/currency/usd/range?from=2024-10-11&to=2024-10-28' \
--header 'Authorization: Bearer Token'
```

**Responses:**
- **200 OK**
  ```json
  {
      "data": {
          "rates": [
              {
                  "id": 5,
                  "created_at": "2024-10-28T11:10:15.406158Z",
                  "updated_at": "2024-10-28T11:10:15.406158Z",
                  "name": "usd",
                  "date": "2024-10-28T00:00:00Z",
                  "rate": 0.011285624
              },
              {
                  "id": 1,
                  "created_at": "2024-10-27T21:34:11.186885Z",
                  "updated_at": "2024-10-27T21:34:11.186885Z",
                  "name": "usd",
                  "date": "2024-10-27T00:00:00Z",
                  "rate": 0.010286624
              }
          ]
      }
  }
  ```
- **400 Bad Request**
  ```json
  {
      "error": "validation error: correct from and to are required"
  }
  ```
- **401 Unauthorized**
  ```json
  {
      "error": "token expired or invalid"
  }
  ```
- **500 Internal Server Error**
  ```json
  {
      "error": "something went wrong"
  }
  ```

The gateway API responses include public error messages (especially for 500 errors), but all errors occurring within the services are logged.