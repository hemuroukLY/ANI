package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/kubercloud/ani/pkg/bootstrap"
	sharedrepo "github.com/kubercloud/ani/pkg/repo"
	"github.com/kubercloud/ani/services/task-service/internal/config"
	"github.com/kubercloud/ani/services/task-service/internal/service"
	"github.com/kubercloud/ani/services/task-service/internal/worker"
)

func main() {
	cfg := config.Load()
	deps := bootstrap.MustConnect(cfg.Config)
	defer deps.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	taskRepo := sharedrepo.NewPostgresAsyncTaskRepo()
	outboxRepo := sharedrepo.NewPostgresOutboxRepo()

	if cfg.OutboxEnabled {
		publisher := worker.NewOutboxPublisher(deps.DB, deps.Ports.MessageBus, outboxRepo, worker.OutboxPublisherConfig{
			PollInterval: cfg.OutboxPollInterval,
			BatchSize:    cfg.OutboxBatchSize,
		}, deps.Logger)
		go publisher.Run(ctx)
	}

	taskSvc := service.NewTaskService(deps.DB, taskRepo)
	bootstrap.RunGRPC(cfg.GRPCPort, taskSvc.Register, deps)
}
