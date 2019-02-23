package crawler

import (
	"sync"

	"github.com/timtosi/mcrawler/internal/domain"
)

// Pipe is an `interface` used by `*crawler.Crawler` to build a crawling
// pipeline.
type Pipe interface {
	Pipe(*sync.WaitGroup, <-chan *domain.Target, chan<- *domain.Target)
}
