package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	taskv1 "github.com/kubercloud/ani/pkg/generated/pb/task/v1"
	sharedrepo "github.com/kubercloud/ani/pkg/repo"
	"github.com/kubercloud/ani/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskService struct {
	taskv1.UnimplementedTaskServiceServer
	db   *pgxpool.Pool
	repo sharedrepo.AsyncTaskRepo
}

func NewTaskService(db *pgxpool.Pool, taskRepo sharedrepo.AsyncTaskRepo) *TaskService {
	return &TaskService{db: db, repo: taskRepo}
}

func (s *TaskService) Register(server *grpc.Server) {
	taskv1.RegisterTaskServiceServer(server, s)
}

func (s *TaskService) GetTask(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.AsyncTask, error) {
	tenantID, taskID, err := parseTenantAndID(req.GetTenantId(), req.GetTaskId())
	if err != nil {
		return nil, toStatus(err)
	}
	ctx = types.WithTenant(ctx, &types.TenantContext{TenantID: tenantID})
	task, err := s.repo.GetByID(ctx, s.db, taskID)
	if err != nil {
		return nil, toStatus(err)
	}
	return taskToPB(task), nil
}

func (s *TaskService) CancelTask(ctx context.Context, req *taskv1.CancelTaskRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "task cancellation requires explicit cancel state transition design")
}

func (s *TaskService) UpdateTaskProgress(ctx context.Context, req *taskv1.UpdateTaskProgressRequest) (*emptypb.Empty, error) {
	tenantID, taskID, workerID, err := parseWorkerMutation(req.GetTenantId(), req.GetTaskId(), req.GetWorkerId())
	if err != nil {
		return nil, toStatus(err)
	}
	ctx = types.WithTenant(ctx, &types.TenantContext{TenantID: tenantID})
	if err := s.repo.UpdateProgress(ctx, s.db, taskID, workerID, int(req.GetProgress())); err != nil {
		return nil, toStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *TaskService) AcquireTaskLease(ctx context.Context, req *taskv1.AcquireTaskLeaseRequest) (*taskv1.AcquireTaskLeaseResponse, error) {
	tenantID, taskID, workerID, err := parseWorkerMutation(req.GetTenantId(), req.GetTaskId(), req.GetWorkerId())
	if err != nil {
		return nil, toStatus(err)
	}
	duration, err := parseLeaseDuration(req.GetLeaseSeconds())
	if err != nil {
		return nil, toStatus(err)
	}
	ctx = types.WithTenant(ctx, &types.TenantContext{TenantID: tenantID})
	acquired, leaseUntil, err := s.repo.AcquireLease(ctx, s.db, taskID, workerID, duration)
	if err != nil {
		return nil, toStatus(err)
	}
	if !acquired {
		return nil, status.Error(codes.AlreadyExists, "task lease already held")
	}
	return &taskv1.AcquireTaskLeaseResponse{
		Acquired:   true,
		LeaseUntil: timestamppb.New(leaseUntil),
	}, nil
}

func (s *TaskService) HeartbeatTaskLease(ctx context.Context, req *taskv1.HeartbeatTaskLeaseRequest) (*emptypb.Empty, error) {
	tenantID, taskID, workerID, err := parseWorkerMutation(req.GetTenantId(), req.GetTaskId(), req.GetWorkerId())
	if err != nil {
		return nil, toStatus(err)
	}
	duration, err := parseLeaseDuration(req.GetLeaseSeconds())
	if err != nil {
		return nil, toStatus(err)
	}
	ctx = types.WithTenant(ctx, &types.TenantContext{TenantID: tenantID})
	if err := s.repo.Heartbeat(ctx, s.db, taskID, workerID, duration); err != nil {
		return nil, toStatus(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *TaskService) FailTask(ctx context.Context, req *taskv1.FailTaskRequest) (*emptypb.Empty, error) {
	tenantID, taskID, workerID, err := parseWorkerMutation(req.GetTenantId(), req.GetTaskId(), req.GetWorkerId())
	if err != nil {
		return nil, toStatus(err)
	}
	if req.GetErrorMessage() == "" {
		return nil, toStatus(types.Wrapf(types.ErrBadRequest, "error_message required"))
	}
	ctx = types.WithTenant(ctx, &types.TenantContext{TenantID: tenantID})
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	defer rollback(ctx, tx)
	if err := s.repo.Fail(ctx, tx, taskID, workerID, req.GetErrorMessage(), req.GetCompensatingAction()); err != nil {
		return nil, toStatus(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &emptypb.Empty{}, nil
}

func (s *TaskService) CompleteTask(ctx context.Context, req *taskv1.CompleteTaskRequest) (*emptypb.Empty, error) {
	tenantID, taskID, workerID, err := parseWorkerMutation(req.GetTenantId(), req.GetTaskId(), req.GetWorkerId())
	if err != nil {
		return nil, toStatus(err)
	}
	var result any = map[string]any{}
	if len(req.GetResult()) > 0 {
		if err := json.Unmarshal(req.GetResult(), &result); err != nil {
			return nil, toStatus(types.Wrapf(types.ErrBadRequest, "result must be valid JSON"))
		}
	}
	ctx = types.WithTenant(ctx, &types.TenantContext{TenantID: tenantID})
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	defer rollback(ctx, tx)
	if err := s.repo.Complete(ctx, tx, taskID, workerID, result); err != nil {
		return nil, toStatus(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &emptypb.Empty{}, nil
}

func parseTenantAndID(tenantID, id string) (uuid.UUID, uuid.UUID, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil || tid == uuid.Nil {
		return uuid.Nil, uuid.Nil, types.Wrapf(types.ErrBadRequest, "invalid tenant_id")
	}
	parsedID, err := uuid.Parse(id)
	if err != nil || parsedID == uuid.Nil {
		return uuid.Nil, uuid.Nil, types.Wrapf(types.ErrBadRequest, "invalid task_id")
	}
	return tid, parsedID, nil
}

func parseWorkerMutation(tenantID, taskID, workerID string) (uuid.UUID, uuid.UUID, string, error) {
	tid, parsedTaskID, err := parseTenantAndID(tenantID, taskID)
	if err != nil {
		return uuid.Nil, uuid.Nil, "", err
	}
	if workerID == "" {
		return uuid.Nil, uuid.Nil, "", types.Wrapf(types.ErrBadRequest, "worker_id required")
	}
	return tid, parsedTaskID, workerID, nil
}

func parseLeaseDuration(seconds int32) (time.Duration, error) {
	if seconds <= 0 || seconds > 3600 {
		return 0, types.Wrapf(types.ErrBadRequest, "lease_seconds must be between 1 and 3600")
	}
	return time.Duration(seconds) * time.Second, nil
}

func rollback(ctx context.Context, tx pgx.Tx) {
	_ = tx.Rollback(ctx)
}

func taskToPB(task *sharedrepo.AsyncTask) *taskv1.AsyncTask {
	out := &taskv1.AsyncTask{
		TenantId:           task.TenantID.String(),
		Id:                 task.ID.String(),
		IdempotencyKey:     task.IdempotencyKey,
		TaskType:           task.TaskType,
		ResourceType:       task.ResourceType,
		ResourceId:         task.ResourceID.String(),
		Status:             task.Status,
		AttemptCount:       int32(task.AttemptCount),
		MaxAttempts:        int32(task.MaxAttempts),
		ProgressPct:        int32(task.ProgressPct),
		LeaseOwner:         task.LeaseOwner,
		ErrorMessage:       task.ErrorMessage,
		CompensatingAction: task.CompensatingAction,
		WebhookUrl:         task.WebhookURL,
		CreatedAt:          timestamppb.New(task.CreatedAt),
	}
	if task.LeaseUntil != nil {
		out.LeaseUntil = timestamppb.New(*task.LeaseUntil)
	}
	if task.LastHeartbeatAt != nil {
		out.LastHeartbeatAt = timestamppb.New(*task.LastHeartbeatAt)
	}
	if task.DeadLetterAt != nil {
		out.DeadLetterAt = timestamppb.New(*task.DeadLetterAt)
	}
	if task.StartedAt != nil {
		out.StartedAt = timestamppb.New(*task.StartedAt)
	}
	if task.CompletedAt != nil {
		out.CompletedAt = timestamppb.New(*task.CompletedAt)
	}
	return out
}

func toStatus(err error) error {
	switch {
	case errors.Is(err, types.ErrNotFound):
		return status.Error(codes.NotFound, "task not found")
	case errors.Is(err, types.ErrBadRequest), errors.Is(err, types.ErrInvalidState):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, types.ErrLeaseTaken):
		return status.Error(codes.AlreadyExists, "task lease is not held by this worker")
	case errors.Is(err, types.ErrForbidden):
		return status.Error(codes.PermissionDenied, "forbidden")
	case errors.Is(err, types.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, "unauthorized")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
