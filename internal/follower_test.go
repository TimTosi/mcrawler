package internal

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/timtosi/mcrawler/internal/domain"
)

func TestFollower_NewFollower(t *testing.T) {
	testCases := []struct {
		name               string
		mockOriginHost     string
		expectOriginHost   string
		expectedAssertFunc func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			"regular",
			"https://www.youtube.com/watch?v=WBupia9oidU",
			"www.youtube.com",
			assert.Nil,
		},
		{
			"noScheme",
			"www.youtube.com/watch?v=WBupia9oidU",
			"",
			assert.NotNil,
		},
		{
			"empty",
			"",
			"",
			assert.NotNil,
		},
		{
			"noHost",
			"/yes",
			"",
			assert.NotNil,
		},
		{
			"noHost_2",
			"https://",
			"",
			assert.NotNil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewFollower(tc.mockOriginHost)
			tc.expectedAssertFunc(t, err)
			if f != nil {
				assert.Equal(t, tc.expectOriginHost, f.originHost)
			}
		})
	}
}

func TestFollower_IsSameHost(t *testing.T) {
	testCases := []struct {
		name               string
		mockOriginHost     string
		mockLink           string
		expectSameHost     bool
		expectedAssertFunc func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			"regular_same",
			"https://www.youtube.com",
			"https://www.youtube.com/watch?v=0NBkIq_xWwI",
			true,
			assert.Nil,
		},
		{
			"regular_same2",
			"http://www.followTest.com",
			"https://www.followTest.com",
			true,
			assert.Nil,
		},
		{
			"noScheme",
			"https://www.testFollow.com",
			"www.testFollow.com/ok?v=test",
			false,
			assert.NotNil,
		},
		{
			"regular_different",
			"https://www.testFollow.com",
			"https://www.followTest.com/ok?v=test",
			false,
			assert.Nil,
		},
		{
			"empty",
			"https://www.youtube.com",
			"",
			false,
			assert.NotNil,
		},
		{
			"noHost",
			"https://www.youtube.com",
			"/yes",
			false,
			assert.NotNil,
		},
		{
			"noHost_2",
			"https://www.youtube.com",
			"https://",
			false,
			assert.NotNil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewFollower(tc.mockOriginHost)
			if err != nil {
				log.Fatalf("%s: %v", tc.name, err)
			}

			ok, err := f.IsSameHost(tc.mockLink)
			tc.expectedAssertFunc(t, err)
			assert.Equal(t, tc.expectSameHost, ok)
		})
	}
}

func TestFollower_Pipe(t *testing.T) {
	testCases := []struct {
		name           string
		mockSameHost   bool
		mockOriginHost string
		mockTarget     *domain.Target
		expectedTarget *domain.Target
	}{
		{
			"regular_same",
			true,
			"https://www.youtube.com",
			&domain.Target{BaseURL: "https://www.youtube.com/watch?v=0NBkIq_xWwI"},
			&domain.Target{BaseURL: "https://www.youtube.com/watch?v=0NBkIq_xWwI"},
		},
		{
			"regular_different",
			false,
			"https://www.youtube.com",
			&domain.Target{BaseURL: "https://www.test.com"},
			nil,
		},
		{
			"empty",
			false,
			"https://www.youtube.com",
			&domain.Target{BaseURL: ""},
			nil,
		},
		{
			"noScheme",
			false,
			"https://www.youtube.com",
			&domain.Target{BaseURL: "www.youtube.com"},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewFollower(tc.mockOriginHost)
			if err != nil {
				log.Fatalf("%s: %v", tc.name, err)
			}

			inChan := make(chan *domain.Target)
			outChan := make(chan *domain.Target)
			wg := sync.WaitGroup{}
			wg.Add(1)

			go f.Pipe(&wg, inChan, outChan)
			inChan <- tc.mockTarget

			select {
			case res := <-outChan:
				assert.Equal(t, tc.expectedTarget, res)
			case <-time.After(1 * time.Second):
				if tc.mockSameHost {
					t.Errorf("%s timeout", tc.name)
				}
			}
		})
	}
}
