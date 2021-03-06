package extractor

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/timtosi/mcrawler/internal/domain"
	"golang.org/x/net/html"
)

// validateScheme returns `true` `scheme` correspond to `http`, `https` or `ftp`
// or `false` otherwise.
func validateScheme(scheme string) bool {
	if len(scheme) >= 4 && scheme[:4] == "http" ||
		len(scheme) == 3 && scheme[:3] == "ftp" {
		return true
	}
	return false
}

// formatLink is an helper function used to format an URL in the proper way.
// It uses `URL` as the URL where the `link` was retrieved from. It returns
// a string representing a valid `URL` or an empty `string` if an `error` occurs
// during parsing.
func formatLink(URL, link string) (string, error) {
	if len(link) == 0 {
		return "", fmt.Errorf("formatLink: link is empty")
	}

	linkURL, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("formatLink: %v found in %s", err, link)
	}

	if validateScheme(linkURL.Scheme) {
		return link, nil
	}

	baseURL, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("formatLink: %v found in %s", err, baseURL)
	} else if len(baseURL.Host) == 0 {
		return "", fmt.Errorf("formatLink: URL %s incomplete", URL)
	}

	if link[0] == '/' || link[len(link)-1] == '/' {
		link = strings.Trim(link, "/")
	}

	return strings.Join([]string{baseURL.Scheme, "://", baseURL.Host, "/", link}, ""), nil
}

// CheckFunc is a named type representing a function that checks if an
// `html.Token` has a link that can be crawled.
type CheckFunc func(html.Token, html.TokenType) string

// Extractor is a `struct` that extracts links found in a web page according to
// the results of its inner `CheckFunc` functions.
type Extractor struct {
	cf []CheckFunc
}

// NewExtractor returns a new `*extractor.Extractor`.
func NewExtractor(checkFuncs ...CheckFunc) *Extractor {
	e := Extractor{cf: make([]CheckFunc, 0)}
	e.cf = append(e.cf, checkFuncs...)
	return &e
}

// ExtractLinks extracts, cleans and returns a `[]string` of links found in
// `content` and matching any `e.cf` function.
func (e *Extractor) ExtractLinks(baseURL string, content []byte) []string {
	var links []string
	rawLink := ""
	uniqueLinks := make(map[string]bool)
	tokenizer := html.NewTokenizer(bytes.NewReader(content))

	for tokenType := tokenizer.Next(); tokenType != html.ErrorToken; tokenType = tokenizer.Next() {
		tkn := tokenizer.Token()
		for _, cf := range e.cf {

			if rawLink = cf(tkn, tokenType); len(rawLink) == 0 {
				continue
			}

			if link, err := formatLink(baseURL, rawLink); err == nil && link != "" {
				uniqueLinks[link] = true
				break
			} else {
				log.Printf("ExtractLinks: %s => %v ", link, err)
			}

		}
	}

	for k := range uniqueLinks {
		links = append(links, k)
	}
	return links
}

// Pipe connects `in` and `out` together. Any `*domain.Target` received from
// `in` will be parsed and extracted links will be sent to `out`.
//
// NOTE: This function will loop over a channel until `in` is closed. After that
// it will close `out`.
func (e *Extractor) Pipe(wg *sync.WaitGroup, in <-chan *domain.Target, out chan<- *domain.Target) {
	defer close(out)

	for t := range in {
		wg.Add(1)
		go func(tgt *domain.Target) {
			links := e.ExtractLinks(tgt.BaseURL, tgt.Content)
			for _, link := range links {
				wg.Add(1)
				go func(l string) { out <- domain.NewTarget(l) }(link)
			}
			wg.Done()
			wg.Done()
		}(t)
	}
}
