package crawler

import (
	"sync"

	"github.com/timtosi/mcrawler/internal"
	"github.com/timtosi/mcrawler/internal/domain"
)

// Crawler is a `struct` that crawls a website.
type Crawler struct {
	urlFrontier chan *domain.Target
}

// NewCrawler returns a new `*crawler.Crawler`.
func NewCrawler() *Crawler {
	return &Crawler{
		urlFrontier: make(chan *domain.Target, 10),
	}
}

// pipeEnd is the function representing the edge of the internal crawling
// pipeline. It cycles new links found during any `crawler.Pipe` to
// `c.urlFrontier`.
//
// NOTE: This function will loop over a channel until `in` is closed.
//
// NOTE: All `crawler.Pipe`s have to `wg.Done()` each time they discard a
// `*domain.Target` from the pipeline.
func (c *Crawler) pipeEnd(in <-chan *domain.Target) {
	for t := range in {
		c.urlFrontier <- t
	}
}

// Run ties all the `crawler.Pipe`s together and initiates the web crawling
// mechanism.
//
// NOTE: It is highly recommended to insert a `internal.Pipe` keeping track of
// already visited web pages in order to avoid looping indefinitely on the same
// links. The `crawler.Archiver` can be used for this goal.
func (c *Crawler) Run(t *domain.Target, pipeline ...internal.Pipe) error {
	wg := sync.WaitGroup{}

	in := c.urlFrontier
	for _, pipe := range pipeline {
		out := make(chan *domain.Target)
		go pipe.Pipe(&wg, in, out)
		in = out
	}
	go c.pipeEnd(in)

	wg.Add(1)
	c.urlFrontier <- t
	wg.Wait()
	close(c.urlFrontier)

	return nil
}
