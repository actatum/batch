package memory

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/actatum/batch/batch"
)

var (
	maxAttempts = 4
)

// repository implements the batch.Repository interface
type repository struct {
	config *batch.Config
	cache  []*batch.Request
}

// NewMemoryRepository returns an object implementing batch.Repository interface
// using memory as a persistence layer
func NewMemoryRepository(c *batch.Config) batch.Repository {
	return &repository{
		config: c,
		cache:  make([]*batch.Request, 0),
	}
}

// Config returns the repositories configuration
func (r *repository) Config() *batch.Config {
	return r.config
}

// Create adds a new request to the cache
func (r *repository) Create(s *batch.Service, req *batch.Request) (*batch.Result, error) {
	var res *batch.Result
	var err error

	if r.isFull() {
		res, err = r.Flush(s)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	r.cache = append(r.cache, req)

	return res, nil
}

// Flush flushes the cache and posts to the configured endpoint
func (r *repository) Flush(s *batch.Service) (*batch.Result, error) {
	res, err := retry(maxAttempts, 2*time.Second, r.post)
	if err != nil {
		s.Logger.Fatal(err.Error())
	}

	s.Logger.Sugar().Infof("batch size: %d, status code: %d, duration: %v", res.Size, res.Code, res.Duration.String())

	return res, nil
}

// isFull returns true if the cache is at max capacity and false otherwise
func (r *repository) isFull() bool {
	return len(r.cache) == r.config.Size
}

// post makes an http post request to the specified endpoint
// sending the contents of the cache as the request body
func (r *repository) post() (*batch.Result, error) {
	size := len(r.cache)
	b, err := json.Marshal(r.cache)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	res, err := http.Post(r.config.Endpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	end := time.Now()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("post failed after retrying " + strconv.Itoa(maxAttempts-1) + " times, status code: " + strconv.Itoa(res.StatusCode))
	}

	r.cache = make([]*batch.Request, 0)

	return &batch.Result{
		Size:     size,
		Code:     res.StatusCode,
		Duration: end.Sub(start),
	}, nil
}

// retry retries the given function the specified amount of times and after the given sleep duration
func retry(attempts int, sleep time.Duration, fn func() (*batch.Result, error)) (*batch.Result, error) {
	res, err := fn()
	if err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, sleep, fn)
		}

		return nil, err
	}

	return res, nil
}
