package crawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timtosi/mcrawler/internal"
	"github.com/timtosi/mcrawler/internal/domain"
	"github.com/timtosi/mcrawler/internal/extractor"
	"github.com/timtosi/mcrawler/internal/mapper"
)

// mockServer is an helper function only used for test purposes. It returns a
// mock `*httptest.Server` webserver or an `error` if something bad occurs.
func mockServer() (*httptest.Server, error) {
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		return nil, fmt.Errorf("mockServer: %v", err)
	}

	ms := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqURL := r.URL.String()
			switch reqURL {
			case "/home":
				w.WriteHeader(http.StatusOK)
				content, err := ioutil.ReadFile("testdata/home.html")
				if err != nil {
					log.Fatalf("Crawler Run: %v", err)
				}
				w.Write(content)
			case "/about":
				w.WriteHeader(http.StatusOK)
				content, err := ioutil.ReadFile("testdata/about.html")
				if err != nil {
					log.Fatalf("Crawler Run: %v", err)
				}
				w.Write(content)
			case "/team":
				w.WriteHeader(http.StatusOK)
				content, err := ioutil.ReadFile("testdata/team.html")
				if err != nil {
					log.Fatalf("Crawler Run: %v", err)
				}
				w.Write(content)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)
	if err := ms.Listener.Close(); err != nil {
		return nil, fmt.Errorf("mockServer: %v", err)
	}
	ms.Listener = l
	ms.Start()
	return ms, nil
}

func TestCrawler_NewCrawler(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{"regular"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() { NewCrawler() })
		})
	}
}

func TestCrawler_Run(t *testing.T) {
	fs, err := mockServer()
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	tgt := domain.NewTarget("http://localhost:8080/home")
	m := mapper.NewMapper()
	f, err := internal.NewFollower(tgt.BaseURL)
	if err != nil {
		log.Fatalf("TestCrawler_Run: %v", err)
	}

	if err := NewCrawler().Run(
		tgt,
		internal.NewArchiver(),
		m,
		f,
		internal.NewWorker(),
		extractor.NewExtractor(extractor.GetImg, extractor.GetLinkNoFollow),
	); err != nil {
		log.Fatal(err)
	}

	assert.ElementsMatch(
		t,
		[]string{
			"http://localhost:8080/home",
			"http://localhost:8080/team",
			"http://localhost:8080/about",
			"http://localhost:8080/notfound",
			"http://localhost:8080/img1.jpg",
			"http://localhost:8080/img2.png",
			"http://localhost:8080/img3.png",
			"http://external1.com",
			"http://external2.com",
			"ftp://www.run-test.com",
		},
		m.SiteMap(),
	)
}
