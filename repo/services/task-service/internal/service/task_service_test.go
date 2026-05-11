package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kubercloud/ani/pkg/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestParseWorkerMutationRequiresTenantTaskAndWorker(t *testing.T) {
	tenantID := uuid.NewString()
	taskID := uuid.NewString()

	gotTenantID, gotTaskID, gotWorkerID, err := parseWorkerMutation(tenantID, taskID, "worker-a")
	if err != nil {
		t.Fatalf("parseWorkerMutation returned error: %v", err)
	}
	if gotTenantID.String() != tenantID {
		t.Fatalf("tenant id mismatch: got %s want %s", gotTenantID, tenantID)
	}
	if gotTaskID.String() != taskID {
		t.Fatalf("task id mismatch: got %s want %s", gotTaskID, taskID)
	}
	if gotWorkerID != "worker-a" {
		t.Fatalf("worker id mismatch: got %s", gotWorkerID)
	}

	_, _, _, err = parseWorkerMutation(tenantID, taskID, "")
	if !errors.Is(err, types.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest for empty worker_id, got %v", err)
	}
}

func TestParseLeaseDurationBounds(t *testing.T) {
	if _, err := parseLeaseDuration(0); !errors.Is(err, types.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest for zero lease, got %v", err)
	}
	if _, err := parseLeaseDuration(3601); !errors.Is(err, types.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest for excessive lease, got %v", err)
	}
	if got, err := parseLeaseDuration(30); err != nil || got.String() != "30s" {
		t.Fatalf("parseLeaseDuration(30) = %v, %v; want 30s, nil", got, err)
	}
}

func TestToStatusMapsLeaseTakenToAlreadyExists(t *testing.T) {
	err := toStatus(types.Wrapf(types.ErrLeaseTaken, "lease owner mismatch"))
	if status.Code(err) != codes.AlreadyExists {
		t.Fatalf("status code = %s, want %s", status.Code(err), codes.AlreadyExists)
	}
}
