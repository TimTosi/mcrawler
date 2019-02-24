package crawler

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/timtosi/mcrawler/internal/domain"
)

// Follower is a `struct` controlling that the pages crawled are only located
// on the `originHost` host.
type Follower struct {
	originHost string
}

// NewFollower returns a new `*crawler.Follower` or an `error` if something
// bad occurs.
func NewFollower(originURL string) (*Follower, error) {
	urlStruct, err := url.Parse(originURL)
	if err != nil {
		return nil, fmt.Errorf("NewFollower: %v", err)
	} else if len(urlStruct.Host) == 0 {
		return nil, fmt.Errorf("NewFollower: no host found in %s", originURL)
	}
	return &Follower{originHost: urlStruct.Host}, nil
}

// IsSameHost ensures that `link` is located on `originHost`. It returns `true`
// if it is the case, `false` otherwise or an `error` if `link` cannot
// be parsed properly by `url.Parse`.
func (f *Follower) IsSameHost(link string) (bool, error) {
	urlStruct, err := url.Parse(link)
	if err != nil {
		return false, fmt.Errorf("IsSameHost: %v", err)
	} else if len(urlStruct.Host) == 0 {
		return false, fmt.Errorf("IsSameHost: no host found in %s", link)
	}

	if urlStruct.Host != f.originHost {
		return false, nil
	}
	return true, nil
}

// Pipe connects `in` and `out` together. Any `*domain.Target` received from
// `in` will be checked against `f.originHost` and be discarded if its host does
// not match `f.originHost`.
//
// NOTE: This function will loop over a channel until `in` is closed. After that
// it will close `out`.
func (f *Follower) Pipe(wg *sync.WaitGroup, in <-chan *domain.Target, out chan<- *domain.Target) {
	defer close(out)

	for t := range in {
		if ok, err := f.IsSameHost(t.BaseURL); err != nil {
			log.Printf("Follower: %f", err)
			wg.Done()
		} else if !ok {
			wg.Done()
		} else {
			out <- t
		}
	}
}
