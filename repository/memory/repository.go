package memory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/actatum/batch/batch"
	errs "github.com/pkg/errors"
)

var (
	maxRetries = 3
)

type repository struct {
	config *batch.Config
	index  int
	cache  []*batch.Request
}

// NewMemoryRepository returns an object implementing batch.Repository interface
// using memory as a persistence layer
func NewMemoryRepository(c *batch.Config) batch.Repository {
	return &repository{
		config: c,
		index:  0,
		cache:  make([]*batch.Request, c.Size),
	}
}

// Create adds a new request to the cache
func (r *repository) Create(req *batch.Request) {
	fmt.Println(len(r.cache))
	if r.isFull() {
		fmt.Println("cache is full")
		if err := r.flush(); err != nil {
			log.Println(err)
		}
		return
	}
	r.cache[r.index] = req
	fmt.Println(r.cache)
	fmt.Println(r.cache[0])
	r.index++
}

func (r *repository) isFull() bool {
	if r.cache[len(r.cache)-1] != nil {
		return true
	}

	return false
}

func (r *repository) flush() error {
	err := retry(2*time.Second, func() (err error) {
		err = r.post()
		return
	})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *repository) post() error {
	b, err := json.Marshal(r.cache)
	if err != nil {
		return errs.Wrap(err, "repository.Memory.flush")
	}

	res, err := http.Post(r.config.Endpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return errs.Wrap(err, "repository.Memory.flush")
	}
	defer res.Body.Close()

	r.cache = make([]*batch.Request, r.config.Size)
	fmt.Println(r.cache)

	return nil
}

func retry(sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}

		if i >= (maxRetries) {
			break
		}

		time.Sleep(sleep)

		log.Println("retrying after error: ", err)
	}
	// TODO: Should be a fatal log here
	return fmt.Errorf("after %d attempts, last error: %s", maxRetries, err)
}
