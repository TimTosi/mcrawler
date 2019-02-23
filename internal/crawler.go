package crawler

import (
	"sync"

	"github.com/timtosi/mcrawler/internal/domain"
)

// Crawler is a `struct` that crawls a website.
type Crawler struct {
	urlFrontier chan *domain.Target
}

// NewCrawler returns a new `*crawler.Crawler`.
func NewCrawler() *Crawler {
	return &Crawler{
		urlFrontier: make(chan *domain.Target),
	}
}

// pipeEnd is the function representing the edge of the internal crawling
// pipeline. It cycles new links found during any `crawler.Pipe` to
// `c.urlFrontier` and manages the `sync.WaitGroup` coordinating the pipeline.
//
// NOTE: This function will loop over a channel until `in` is closed.
//
// NOTE: `crawler.Pipe`s have to `wg.Done()` each time they discard a
// `*domain.Target` from the pipeline.
func (c *Crawler) pipeEnd(wg *sync.WaitGroup, in <-chan *domain.Target) {
	for t := range in {
		if !t.Done {
			wg.Add(1)
			c.urlFrontier <- t
		} else {
			wg.Done()
		}
	}
}

// Run ties all the `crawler.Pipe`s together and initiates the web crawling
// mechanism.
//
// NOTE: It is highly recommended to insert a `crawler.Pipe` keeping track of
// already visited web pages in order to avoid looping indefinitely on the same
// links. The `crawler.Archiver` can be used for this goal.
func (c *Crawler) Run(t *domain.Target, pipeline ...Pipe) error {
	wg := sync.WaitGroup{}

	wg.Add(1)
	c.urlFrontier <- t

	in := c.urlFrontier
	for _, pipe := range pipeline {
		out := make(chan *domain.Target)
		go pipe.Pipe(&wg, in, out)
		in = out
	}
	go c.pipeEnd(&wg, in)

	wg.Wait()
	close(c.urlFrontier)

	return nil
}
