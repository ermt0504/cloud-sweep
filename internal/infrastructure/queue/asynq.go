package queue

import (
	"github.com/cloudsweep/cloudsweep/internal/infrastructure/config"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// Task types
const (
	TaskTypeScanResources    = "scan:resources"
	TaskTypeCleanupResources = "cleanup:resources"
	TaskTypeApplyPolicy      = "policy:apply"
	TaskTypeSendNotification = "notification:send"
)

// NewAsynqClient creates a new Asynq client
func NewAsynqClient(cfg config.RedisConfig) (*asynq.Client, error) {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return client, nil
}

// NewWorkerServer creates a new Asynq server for processing tasks
func NewWorkerServer(cfg config.RedisConfig, db *gorm.DB) (*asynq.Server, error) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return srv, nil
}

// NewServeMux creates a new Asynq ServeMux with handlers
func NewServeMux(db *gorm.DB) *asynq.ServeMux {
	mux := asynq.NewServeMux()

	// Register handlers
	mux.HandleFunc(TaskTypeScanResources, HandleScanResources(db))
	mux.HandleFunc(TaskTypeCleanupResources, HandleCleanupResources(db))
	mux.HandleFunc(TaskTypeApplyPolicy, HandleApplyPolicy(db))
	mux.HandleFunc(TaskTypeSendNotification, HandleSendNotification(db))

	return mux
}
