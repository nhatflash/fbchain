package tasks

import (
	"context"
	"time"
	"os"
	"github.com/hibiken/asynq"
	"github.com/nhatflash/fbchain/repository"
)

const (
	CleanUpExpiredOrderTemplate string = "expiredOrders:cleanUp"	
)

func RegisterAsynqClient() *asynq.Client {
	addr := os.Getenv("REDIS_SERVER")
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr: addr,
	})
}


func NewCleanUpExpiredOrdersTask() *asynq.Task {
	return asynq.NewTask(CleanUpExpiredOrderTemplate, nil, asynq.Unique(5*time.Minute))
}


type CleanUpOrdersTaskHandler struct {
	ROrderRepo 			*repository.RestaurantOrderRepository
}

func (o *CleanUpOrdersTaskHandler) HandleCleanUpExpiredOrdersTask(ctx context.Context, t *asynq.Task) error {
	err := o.ROrderRepo.DeleteExpiredPendingOrders(ctx)
	if err != nil {
		return err
	}
	return nil
}


func RegisterAsynqServer(ror *repository.RestaurantOrderRepository) (*asynq.Server, *asynq.ServeMux) {
	addr := os.Getenv("REDIS_SERVER")

	// calling handlers
	cleanUpOrdersHandler := &CleanUpOrdersTaskHandler{
		ROrderRepo: ror,
	}

	// initialize asynq server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: addr},
		asynq.Config{
			Concurrency: 10,
			ShutdownTimeout: 8*time.Second,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(CleanUpExpiredOrderTemplate, cleanUpOrdersHandler.HandleCleanUpExpiredOrdersTask)

	return srv, mux
}


func RegisterAsynqScheduler() (*asynq.Scheduler, error) {
	addr := os.Getenv("REDIS_SERVER")
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: addr},
		&asynq.SchedulerOpts{},
	)

	var cleanUpOrdersTask *asynq.Task
	cleanUpOrdersTask = NewCleanUpExpiredOrdersTask()

	var err error
	_, err = scheduler.Register("@every 5m", cleanUpOrdersTask)
	if err != nil {
		return nil, err
	}
	return scheduler, nil
}