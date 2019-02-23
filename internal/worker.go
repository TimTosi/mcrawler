package crawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/timtosi/mcrawler/internal/domain"
)

// Worker is a `struct` representing a HTTP client concurrently fetching
// web pages.
type Worker struct {
	http.Client
}

// NewWorker returns a new `*crawler.Worker` that can be configured
// through `opts` functions.
func NewWorker(opts ...func(*Worker)) *Worker {
	tr := &http.Transport{
		DisableKeepAlives:     true,
		MaxIdleConnsPerHost:   1,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}

	w := &Worker{Client: http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}}

	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Fetch performs a `GET` request on the web page located at `t.BaseURL` and
// populates its `t.Content` or returns an `error` if something bad occurs.
func (w *Worker) Fetch(t *domain.Target) error {
	resp, err := w.Get(t.BaseURL)
	if err != nil {
		return fmt.Errorf("Fetch: %v", err)
	}

	t.Content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Fetch: %v", err)
	}

	if err = resp.Body.Close(); err != nil {
		return fmt.Errorf("Fetch: %v", err)
	}
	return nil
}

// Pipe connects `in` and `out` together. Any `*domain.Target` received from
// `in` will be fetched and the web page content will be sent to `out` if no
// error occurs.
//
// NOTE: This function will loop over a channel until `in` is closed. After that
// it will close `out`.
func (w *Worker) Pipe(wg *sync.WaitGroup, in <-chan *domain.Target, out chan<- *domain.Target) {
	defer close(out)

	for t := range in {
		if err := w.Fetch(t); err != nil {
			log.Printf("Worker: %f", err)
			wg.Done()
		} else {
			t.Done = true
			out <- t
		}
	}
}
