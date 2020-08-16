package transport

import (
	"os"
	"strconv"

	"github.com/actatum/batch/batch"
	"github.com/actatum/batch/repository/memory"
	errs "github.com/pkg/errors"
)

// Run starts the http server
func Run() error {
	conf, err := getConfig()
	if err != nil {
		return errs.Wrap(err, "transport.Run")
	}

	repo := memory.NewMemoryRepository(conf)

	logger, err := batch.NewLogger()
	if err != nil {
		return errs.Wrap(err, "transport.Run")
	}

	service := batch.NewBatchService(repo, logger)

	server := NewServer(service)

	r := routes(server)

	go service.Background()

	service.Logger.Info("Batch server initialized at port :8080")

	return r.Start(":8080")
}

// getConfig retrieves necessary configuration from the environment
func getConfig() (*batch.Config, error) {
	s := os.Getenv("BATCH_SIZE")
	i := os.Getenv("BATCH_INTERVAL")
	e := os.Getenv("BATCH_ENDPOINT")

	size, err := strconv.Atoi(s)
	if err != nil {
		return nil, errs.Wrap(err, "transport.getConfig")
	}

	interval, err := strconv.Atoi(i)
	if err != nil {
		return nil, errs.Wrap(err, "transport.getConfig")
	}

	return &batch.Config{
		Size:     size,
		Interval: interval,
		Endpoint: e,
	}, nil
}
