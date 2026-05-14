// Package repo contains shared repository implementations used by ANI services.
package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	natsmsg "github.com/kubercloud/ani/pkg/nats"
	"github.com/kubercloud/ani/pkg/types"
)

// AsyncTaskRepo stores async task state and the outbox events that drive NATS.
type AsyncTaskRepo interface {
	Create(ctx context.Context, tx pgx.Tx, req CreateTaskReq) (*AsyncTask, error)
	GetByID(ctx context.Context, pool *pgxpool.Pool, id uuid.UUID) (*AsyncTask, error)
	GetByIdempotencyKey(ctx context.Context, pool *pgxpool.Pool, tenantID uuid.UUID, key string) (*AsyncTask, error)
	AcquireLease(ctx context.Context, pool *pgxpool.Pool, taskID uuid.UUID, workerID string, duration time.Duration) (bool, time.Time, error)
	Heartbeat(ctx context.Context, pool *pgxpool.Pool, taskID uuid.UUID, workerID string, duration time.Duration) error
	UpdateProgress(ctx context.Context, pool *pgxpool.Pool, taskID uuid.UUID, workerID string, pct int) error
	Complete(ctx context.Context, tx pgx.Tx, taskID uuid.UUID, workerID string, result any) error
	Fail(ctx context.Context, tx pgx.Tx, taskID uuid.UUID, workerID string, errMsg, compensatingAction string) error
	GetExpiredLeases(ctx context.Context, pool *pgxpool.Pool, limit int) ([]*AsyncTask, error)
}

// PostgresAsyncTaskRepo is the PostgreSQL implementation of AsyncTaskRepo.
type PostgresAsyncTaskRepo struct{}

func NewPostgresAsyncTaskRepo() *PostgresAsyncTaskRepo {
	return &PostgresAsyncTaskRepo{}
}

type CreateTaskReq struct {
	TenantID       uuid.UUID
	IdempotencyKey string
	TaskType       string
	ResourceType   string
	ResourceID     uuid.UUID
	MaxAttempts    int
	WebhookURL     string
	OutboxSubject  string
	OutboxPayload  any
}

type AsyncTask struct {
	TenantID           uuid.UUID
	ID                 uuid.UUID
	IdempotencyKey     string
	TaskType           string
	ResourceType       string
	ResourceID         uuid.UUID
	Status             string
	AttemptCount       int
	MaxAttempts        int
	ProgressPct        int
	LeaseOwner         string
	ErrorMessage       string
	CompensatingAction string
	LeaseUntil         *time.Time
	LastHeartbeatAt    *time.Time
	DeadLetterAt       *time.Time
	WebhookURL         string
	CreatedAt          time.Time
	StartedAt          *time.Time
	CompletedAt        *time.Time
}

func (r *PostgresAsyncTaskRepo) Create(ctx context.Context, tx pgx.Tx, req CreateTaskReq) (*AsyncTask, error) {
	if err := validateCreateTaskReq(req); err != nil {
		return nil, err
	}
	if err := types.SetDBTenant(ctx, tx); err != nil {
		return nil, fmt.Errorf("taskRepo.Create set tenant: %w", err)
	}

	payload, err := json.Marshal(req.OutboxPayload)
	if err != nil {
		return nil, fmt.Errorf("taskRepo.Create marshal outbox payload: %w", err)
	}
	if req.MaxAttempts <= 0 {
		req.MaxAttempts = 3
	}

	var task AsyncTask
	err = tx.QueryRow(ctx, `
		INSERT INTO async_tasks (
			tenant_id, idempotency_key, task_type, resource_type, resource_id,
			max_attempts, webhook_url, payload
		)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, '00000000-0000-0000-0000-000000000000')::uuid,
			$6, NULLIF($7, ''), $8)
		ON CONFLICT (tenant_id, idempotency_key) DO NOTHING
		RETURNING tenant_id, id, idempotency_key, task_type, COALESCE(resource_type, ''),
			COALESCE(resource_id, '00000000-0000-0000-0000-000000000000'::uuid),
			status, attempt_count, max_attempts, progress_pct,
			COALESCE(lease_owner, ''),
			COALESCE(error_message, ''), COALESCE(compensating_action, ''),
			lease_until, last_heartbeat_at, dead_letter_at, COALESCE(webhook_url, ''),
			created_at, started_at, completed_at
	`, req.TenantID, req.IdempotencyKey, req.TaskType, req.ResourceType, req.ResourceID,
		req.MaxAttempts, req.WebhookURL, payload).Scan(taskScanDest(&task)...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, types.Wrapf(types.ErrConflict, "taskRepo.Create idempotency_key=%s", req.IdempotencyKey)
	}
	if err != nil {
		return nil, fmt.Errorf("taskRepo.Create insert task: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO outbox_events (
			aggregate_type, aggregate_id, event_type, tenant_id, payload
		)
		VALUES ($1, $2, $3, $4, $5)
	`, req.ResourceType, task.ID, req.OutboxSubject, req.TenantID, payload); err != nil {
		return nil, fmt.Errorf("taskRepo.Create insert outbox: %w", err)
	}

	return &task, nil
}

func (r *PostgresAsyncTaskRepo) GetByID(ctx context.Context, pool *pgxpool.Pool, id uuid.UUID) (*AsyncTask, error) {
	tx, err := beginTenantTx(ctx, pool)
	if err != nil {
		return nil, err
	}
	defer rollback(ctx, tx)

	task, err := getTaskByQuery(ctx, tx, `WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("taskRepo.GetByID commit: %w", err)
	}
	return task, nil
}

func (r *PostgresAsyncTaskRepo) GetByIdempotencyKey(ctx context.Context, pool *pgxpool.Pool, tenantID uuid.UUID, key string) (*AsyncTask, error) {
	tx, err := beginTenantTx(ctx, pool)
	if err != nil {
		return nil, err
	}
	defer rollback(ctx, tx)

	task, err := getTaskByQuery(ctx, tx, `WHERE tenant_id=$1 AND idempotency_key=$2`, tenantID, key)
	if errors.Is(err, types.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("taskRepo.GetByIdempotencyKey commit: %w", err)
	}
	return task, nil
}

func (r *PostgresAsyncTaskRepo) AcquireLease(ctx context.Context, pool *pgxpool.Pool, taskID uuid.UUID, workerID string, duration time.Duration) (bool, time.Time, error) {
	if workerID == "" {
		return false, time.Time{}, types.Wrapf(types.ErrBadRequest, "taskRepo.AcquireLease worker_id required")
	}
	if duration <= 0 {
		return false, time.Time{}, types.Wrapf(types.ErrBadRequest, "taskRepo.AcquireLease duration must be positive")
	}
	tx, err := beginTenantTx(ctx, pool)
	if err != nil {
		return false, time.Time{}, err
	}
	defer rollback(ctx, tx)

	var leaseUntil time.Time
	err = tx.QueryRow(ctx, `
		UPDATE async_tasks
		SET status='running',
			started_at=COALESCE(started_at, NOW()),
			lease_owner=$2,
			lease_until=NOW() + ($3::double precision * INTERVAL '1 second'),
			last_heartbeat_at=NOW(),
			updated_at=NOW()
		WHERE id=$1
		  AND status IN ('pending', 'running', 'failed')
		  AND (lease_until IS NULL OR lease_until < NOW())
		RETURNING lease_until
	`, taskID, workerID, duration.Seconds()).Scan(&leaseUntil)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, time.Time{}, nil
	}
	if err != nil {
		return false, time.Time{}, fmt.Errorf("taskRepo.AcquireLease update: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return false, time.Time{}, fmt.Errorf("taskRepo.AcquireLease commit: %w", err)
	}
	return true, leaseUntil, nil
}

func (r *PostgresAsyncTaskRepo) Heartbeat(ctx context.Context, pool *pgxpool.Pool, taskID uuid.UUID, workerID string, duration time.Duration) error {
	if workerID == "" {
		return types.Wrapf(types.ErrBadRequest, "taskRepo.Heartbeat worker_id required")
	}
	if duration <= 0 {
		return types.Wrapf(types.ErrBadRequest, "taskRepo.Heartbeat duration must be positive")
	}
	tx, err := beginTenantTx(ctx, pool)
	if err != nil {
		return err
	}
	defer rollback(ctx, tx)

	tag, err := tx.Exec(ctx, `
		UPDATE async_tasks
		SET last_heartbeat_at=NOW(),
			lease_until=NOW() + ($3::double precision * INTERVAL '1 second'),
			updated_at=NOW()
		WHERE id=$1 AND status='running' AND lease_owner=$2 AND lease_until >= NOW()
	`, taskID, workerID, duration.Seconds())
	if err != nil {
		return fmt.Errorf("taskRepo.Heartbeat update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return types.Wrapf(types.ErrLeaseTaken, "taskRepo.Heartbeat task_id=%s worker_id=%s", taskID, workerID)
	}
	return tx.Commit(ctx)
}

func (r *PostgresAsyncTaskRepo) UpdateProgress(ctx context.Context, pool *pgxpool.Pool, taskID uuid.UUID, workerID string, pct int) error {
	if workerID == "" {
		return types.Wrapf(types.ErrBadRequest, "taskRepo.UpdateProgress worker_id required")
	}
	if pct < 0 || pct > 100 {
		return types.Wrapf(types.ErrBadRequest, "taskRepo.UpdateProgress pct=%d", pct)
	}
	tx, err := beginTenantTx(ctx, pool)
	if err != nil {
		return err
	}
	defer rollback(ctx, tx)

	tag, err := tx.Exec(ctx, `
		UPDATE async_tasks
		SET progress_pct=$3, last_heartbeat_at=NOW(), updated_at=NOW()
		WHERE id=$1 AND status='running' AND lease_owner=$2 AND lease_until >= NOW()
	`, taskID, workerID, pct)
	if err != nil {
		return fmt.Errorf("taskRepo.UpdateProgress update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return types.Wrapf(types.ErrLeaseTaken, "taskRepo.UpdateProgress task_id=%s worker_id=%s", taskID, workerID)
	}
	return tx.Commit(ctx)
}

func (r *PostgresAsyncTaskRepo) Complete(ctx context.Context, tx pgx.Tx, taskID uuid.UUID, workerID string, result any) error {
	if workerID == "" {
		return types.Wrapf(types.ErrBadRequest, "taskRepo.Complete worker_id required")
	}
	if err := types.SetDBTenant(ctx, tx); err != nil {
		return fmt.Errorf("taskRepo.Complete set tenant: %w", err)
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("taskRepo.Complete marshal result: %w", err)
	}

	var tenantID uuid.UUID
	var taskType string
	err = tx.QueryRow(ctx, `
		UPDATE async_tasks
		SET status='completed',
			progress_pct=100,
			result=$2,
			lease_owner=NULL,
			lease_until=NULL,
			completed_at=NOW(),
			updated_at=NOW()
		WHERE id=$1 AND status='running' AND lease_owner=$3 AND lease_until >= NOW()
		RETURNING tenant_id, task_type
	`, taskID, resultJSON, workerID).Scan(&tenantID, &taskType)
	if errors.Is(err, pgx.ErrNoRows) {
		return types.Wrapf(types.ErrLeaseTaken, "taskRepo.Complete task_id=%s worker_id=%s", taskID, workerID)
	}
	if err != nil {
		return fmt.Errorf("taskRepo.Complete update: %w", err)
	}

	event := natsmsg.TaskCompletedEvent{
		TaskID:      taskID,
		TenantID:    tenantID,
		TaskType:    taskType,
		Status:      "completed",
		Result:      result,
		PublishedAt: time.Now().UTC(),
	}
	return insertTaskEvent(ctx, tx, taskID, tenantID, event)
}

func (r *PostgresAsyncTaskRepo) Fail(ctx context.Context, tx pgx.Tx, taskID uuid.UUID, workerID string, errMsg, compensatingAction string) error {
	if workerID == "" {
		return types.Wrapf(types.ErrBadRequest, "taskRepo.Fail worker_id required")
	}
	if err := types.SetDBTenant(ctx, tx); err != nil {
		return fmt.Errorf("taskRepo.Fail set tenant: %w", err)
	}

	var tenantID uuid.UUID
	var taskType string
	var status string
	err := tx.QueryRow(ctx, `
		UPDATE async_tasks
		SET attempt_count=attempt_count+1,
			status=CASE WHEN attempt_count+1 >= max_attempts THEN 'dead_letter' ELSE 'failed' END,
			error_message=$2,
			compensating_action=NULLIF($3, ''),
			dead_letter_at=CASE WHEN attempt_count+1 >= max_attempts THEN NOW() ELSE dead_letter_at END,
			lease_owner=NULL,
			lease_until=NULL,
			updated_at=NOW()
		WHERE id=$1 AND status='running' AND lease_owner=$4 AND lease_until >= NOW()
		RETURNING tenant_id, task_type, status
	`, taskID, errMsg, compensatingAction, workerID).Scan(&tenantID, &taskType, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return types.Wrapf(types.ErrLeaseTaken, "taskRepo.Fail task_id=%s worker_id=%s", taskID, workerID)
	}
	if err != nil {
		return fmt.Errorf("taskRepo.Fail update: %w", err)
	}

	event := natsmsg.TaskCompletedEvent{
		TaskID:      taskID,
		TenantID:    tenantID,
		TaskType:    taskType,
		Status:      status,
		ErrorMsg:    errMsg,
		PublishedAt: time.Now().UTC(),
	}
	return insertTaskEvent(ctx, tx, taskID, tenantID, event)
}

func (r *PostgresAsyncTaskRepo) GetExpiredLeases(ctx context.Context, pool *pgxpool.Pool, limit int) ([]*AsyncTask, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	tx, err := beginTenantTx(ctx, pool)
	if err != nil {
		return nil, err
	}
	defer rollback(ctx, tx)

	rows, err := tx.Query(ctx, taskSelectSQL+`
		WHERE status='running' AND lease_until < NOW()
		ORDER BY lease_until ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("taskRepo.GetExpiredLeases query: %w", err)
	}
	defer rows.Close()

	var tasks []*AsyncTask
	for rows.Next() {
		task := &AsyncTask{}
		if err := rows.Scan(taskScanDest(task)...); err != nil {
			return nil, fmt.Errorf("taskRepo.GetExpiredLeases scan: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("taskRepo.GetExpiredLeases rows: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("taskRepo.GetExpiredLeases commit: %w", err)
	}
	return tasks, nil
}

const taskSelectSQL = `
	SELECT tenant_id, id, idempotency_key, task_type, COALESCE(resource_type, ''),
		COALESCE(resource_id, '00000000-0000-0000-0000-000000000000'::uuid),
		status, attempt_count, max_attempts, progress_pct,
		COALESCE(lease_owner, ''),
		COALESCE(error_message, ''), COALESCE(compensating_action, ''),
		lease_until, last_heartbeat_at, dead_letter_at, COALESCE(webhook_url, ''),
		created_at, started_at, completed_at
	FROM async_tasks
`

func getTaskByQuery(ctx context.Context, tx pgx.Tx, where string, args ...any) (*AsyncTask, error) {
	task := &AsyncTask{}
	err := tx.QueryRow(ctx, taskSelectSQL+" "+where, args...).Scan(taskScanDest(task)...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, types.Wrapf(types.ErrNotFound, "taskRepo.Get")
	}
	if err != nil {
		return nil, fmt.Errorf("taskRepo.Get query: %w", err)
	}
	return task, nil
}

func taskScanDest(t *AsyncTask) []any {
	return []any{
		&t.TenantID, &t.ID, &t.IdempotencyKey, &t.TaskType, &t.ResourceType,
		&t.ResourceID, &t.Status, &t.AttemptCount, &t.MaxAttempts, &t.ProgressPct,
		&t.LeaseOwner, &t.ErrorMessage, &t.CompensatingAction, &t.LeaseUntil, &t.LastHeartbeatAt,
		&t.DeadLetterAt, &t.WebhookURL, &t.CreatedAt, &t.StartedAt, &t.CompletedAt,
	}
}

func beginTenantTx(ctx context.Context, pool *pgxpool.Pool) (pgx.Tx, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tenant tx: %w", err)
	}
	if err := types.SetDBTenant(ctx, tx); err != nil {
		_ = tx.Rollback(ctx)
		return nil, fmt.Errorf("set tenant: %w", err)
	}
	return tx, nil
}

func rollback(ctx context.Context, tx pgx.Tx) {
	_ = tx.Rollback(ctx)
}

func insertTaskEvent(ctx context.Context, tx pgx.Tx, taskID, tenantID uuid.UUID, event natsmsg.TaskCompletedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal task event: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO outbox_events (
			aggregate_type, aggregate_id, event_type, tenant_id, payload
		)
		VALUES ('async_task', $1, $2, $3, $4)
	`, taskID, natsmsg.TaskCompletedSubject(taskID), tenantID, payload); err != nil {
		return fmt.Errorf("insert task event outbox: %w", err)
	}
	return nil
}

func validateCreateTaskReq(req CreateTaskReq) error {
	if req.TenantID == uuid.Nil {
		return types.Wrapf(types.ErrBadRequest, "taskRepo.Create tenant_id required")
	}
	if req.IdempotencyKey == "" {
		return types.Wrapf(types.ErrBadRequest, "taskRepo.Create idempotency_key required")
	}
	if req.TaskType == "" {
		return types.Wrapf(types.ErrBadRequest, "taskRepo.Create task_type required")
	}
	if req.OutboxSubject == "" {
		return types.Wrapf(types.ErrBadRequest, "taskRepo.Create outbox subject required")
	}
	return nil
}
