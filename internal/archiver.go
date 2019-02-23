package crawler

import (
	"sync"

	"github.com/timtosi/mcrawler/internal/domain"
)

// Archiver is a `struct` checking that a given URL has not already been seen.
type Archiver struct {
	archive map[string]bool
	mu      *sync.Mutex
}

// NewArchiver returns a new `*crawler.Archiver`
func NewArchiver() *Archiver {
	return &Archiver{
		archive: make(map[string]bool),
		mu:      &sync.Mutex{},
	}
}

// IsAlreadySeen compares `url` to `a.archive` and returns `true` if it is
// already seen. If not, `url` is stored in `a.archive` and the function
// returns `false`.
//
// NOTE: This function is thread-safe.
func (a *Archiver) IsAlreadySeen(url string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.archive[url]; ok {
		return true
	}
	a.archive[url] = true
	return false
}

// Pipe connects `in` and `out` together. Any `t` received from `in` will
// be checked against `a.archive` and sent to `out` if not already seen.
//
// NOTE: This function will loop over a channel until `in` is closed. After that
// it will close `out`.
func (a *Archiver) Pipe(wg *sync.WaitGroup, in <-chan *domain.Target, out chan<- *domain.Target) {
	defer close(out)

	for t := range in {
		if a.IsAlreadySeen(t.BaseURL) {
			wg.Done()
		} else {
			out <- t
		}
	}
}
