package internal

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/timtosi/mcrawler/internal/domain"
)

// mockServer is an helper function only used for test purposes. It returns a
// mock `*httptest.Server` webserver.
func mockServer() *httptest.Server {
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqURL := r.URL.String()
			switch reqURL {
			case "/good":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`correctly retrieved`))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)
	return mockServer
}

// -----------------------------------------------------------------------------

func TestWorker_NewWorker(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{"regular"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() { NewWorker() })
		})
	}
}

func TestWorker_Fetch(t *testing.T) {
	testCases := []struct {
		name               string
		mockURL            string
		expectedContent    string
		expectedAssertFunc func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			"regular",
			"/good",
			"correctly retrieved",
			assert.Nil,
		},
		{
			"pathNotFound",
			"/nope",
			"",
			assert.Nil,
		},
		{
			"badURL",
			"http://",
			"",
			assert.NotNil,
		},
	}

	mockServer := mockServer()
	defer mockServer.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := NewWorker()
			tgt := domain.NewTarget(fmt.Sprintf("%s%s", mockServer.URL, tc.mockURL))

			tc.expectedAssertFunc(t, w.Fetch(tgt))
			assert.Equal(t, tc.expectedContent, string(tgt.Content))
		})
	}
}

func TestWorker_Pipe(t *testing.T) {
	testCases := []struct {
		name            string
		mockURL         string
		expectedContent string
		expectedTimeout bool
	}{
		{
			"regular",
			"/good",
			"correctly retrieved",
			false,
		},
		{
			"pathNotFound",
			"/nope",
			"",
			false,
		},
		{
			"badURL",
			"http://",
			"",
			true,
		},
	}

	ms := mockServer()
	defer ms.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := NewWorker()

			tgt := domain.NewTarget(fmt.Sprintf("%s%s", ms.URL, tc.mockURL))

			inChan := make(chan *domain.Target)
			outChan := make(chan *domain.Target)
			wg := sync.WaitGroup{}
			wg.Add(1)

			go w.Pipe(&wg, inChan, outChan)
			inChan <- tgt

			select {
			case res := <-outChan:
				assert.Equal(t, tc.expectedContent, string(res.Content))
			case <-time.After(1 * time.Second):
				if !tc.expectedTimeout {
					log.Fatalf("%s: test fail", tc.name)
				}
			}
		})
	}
}
