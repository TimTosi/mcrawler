package mapper

import (
	"fmt"
	"sync"

	"github.com/timtosi/mcrawler/internal/domain"
)

// Mapper is a `struct` used for rendering a Site Map.
type Mapper struct {
	siteMap []string
	mu      *sync.RWMutex
}

// NewMapper returns a new `*mapper.Mapper`.
func NewMapper() *Mapper {
	return &Mapper{
		siteMap: make([]string, 0),
		mu:      &sync.RWMutex{},
	}
}

// Add adds `link` to `m.siteMap`.
//
// NOTE: This function is thread-safe.
func (m *Mapper) Add(link string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.siteMap = append(m.siteMap, link)
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

	for _, k := range m.siteMap {
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
		m.Add(t.BaseURL)
		out <- t
	}
}
