package memory

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/actatum/batch/batch"
)

var (
	maxAttempts = 4
	mu          sync.Mutex
)

// repository implements the batch.Repository interface
type repository struct {
	config *batch.Config
	cache  []batch.Request
}

// NewMemoryRepository returns an object implementing batch.Repository interface
// using memory as a persistence layer
func NewMemoryRepository(c *batch.Config) batch.Repository {
	return &repository{
		config: c,
		cache:  make([]batch.Request, 0),
	}
}

// Config returns the repositories configuration
func (r *repository) Config() *batch.Config {
	return r.config
}

// Create adds a new request to the cache
func (r *repository) Add(req batch.Request) {
	mu.Lock()
	r.cache = append(r.cache, req)
	mu.Unlock()
}

// Flush flushes the cache and posts to the configured endpoint
func (r *repository) Flush() (*batch.Result, error) {
	mu.Lock()
	defer mu.Unlock()
	res, err := retry(maxAttempts, 2*time.Second, r.post)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// WillFill returns true if the next entry to the cache will fill it to capacity
func (r *repository) WillFill() bool {
	mu.Lock()
	defer mu.Unlock()
	return len(r.cache)+1 == r.config.Size
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

	r.cache = make([]batch.Request, 0)

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
