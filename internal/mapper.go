package crawler

import (
	"fmt"
	"sync"

	"github.com/timtosi/mcrawler/internal/domain"
)

// Mapper is a `struct` used for rendering a Site Map.
type Mapper struct {
	siteMap map[string]bool
	mu      *sync.RWMutex
}

// NewMapper returns a new `*crawler.Mapper`.
func NewMapper() *Mapper {
	return &Mapper{
		siteMap: make(map[string]bool),
		mu:      &sync.RWMutex{},
	}
}

// Add adds `t.BaseURL` to `m.siteMap`.
//
// NOTE: This function is thread-safe.
func (m *Mapper) Add(t *domain.Target) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.siteMap[t.BaseURL] = true
}

// Render renders a Site Map of urls contained in `m.siteMap` to standard
// output.
//
// NOTE: This function is thread-safe.
func (m *Mapper) Render() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Println(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	for k := range m.siteMap {
		fmt.Println("\t<url>")
		fmt.Printf("\t\t<loc>%s</loc>\n", k)
		fmt.Println("\t</url>")
	}
	fmt.Println("</urlset>")
}

// Pipe connects `in` and `out` together. Any `*domain.Target` received from
// `in` that have been processed will be added to `m.siteMap`.
//
// NOTE: This function will loop over a channel until `in` is closed. After that
// it will close `out`.
func (m *Mapper) Pipe(wg *sync.WaitGroup, in <-chan *domain.Target, out chan<- *domain.Target) {
	defer close(out)

	for t := range in {
		if t.Done {
			m.Add(t)
		}
		out <- t
	}
}
