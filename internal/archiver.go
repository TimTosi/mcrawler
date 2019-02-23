package crawler

import (
	"sync"

	"github.com/timtosi/mcrawler/internal/domain"
)

// Archiver is a `struct` checking that a given URL has not already been seen.
type Archiver struct {
	archive map[string]bool
}

// NewArchiver returns a new `*crawler.Archiver`
func NewArchiver() *Archiver {
	return &Archiver{archive: make(map[string]bool)}
}

// IsAlreadySeen compares `url` to `a.archive` and returns `true` if it is
// already seen. If not, `url` is stored in `a.archive` and the function
// returns `false`.
func (a *Archiver) IsAlreadySeen(url string) bool {
	if _, ok := a.archive[url]; ok {
		return true
	}
	a.archive[url] = true
	return false
}

// GetArchive returns a `[]string` containing every single URL contained in
// `a.archive`.
func (a *Archiver) GetArchive() []string {
	var res []string

	for key := range a.archive {
		res = append(res, key)
	}
	return res
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
