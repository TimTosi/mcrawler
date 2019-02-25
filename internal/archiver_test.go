package crawler

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/timtosi/mcrawler/internal/domain"
)

func TestArchiver_NewArchiver(t *testing.T) {
	testCases := []struct {
		name               string
		expectedAssertFunc func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{"regular", assert.NotNil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectedAssertFunc(t, NewArchiver())
		})
	}
}

func TestArchiver_IsAlreadySeen(t *testing.T) {
	testCases := []struct {
		name             string
		mockArchived     []string
		mockURL          string
		expectedReturn   bool
		expectedArchived []string
	}{
		{
			"empty",
			[]string{},
			"https://www.youtube.com/watch?v=O9FnFxoXBmg&t=1s",
			false,
			[]string{"https://www.youtube.com/watch?v=O9FnFxoXBmg&t=1s"},
		},
		{
			"notSeen",
			[]string{"https://fakeSeen.com"},
			"https://notaSeen.com",
			false,
			[]string{"https://fakeSeen.com", "https://notaSeen.com"},
		},
		{
			"alreadySeen",
			[]string{"https://yesSeen.com"},
			"https://yesSeen.com",
			true,
			[]string{"https://yesSeen.com"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewArchiver()
			res := make([]string, 0)
			for _, k := range tc.mockArchived {
				a.archive[k] = true
			}

			assert.Equal(t, tc.expectedReturn, a.IsAlreadySeen(tc.mockURL))
			for k := range a.archive {
				res = append(res, k)
			}
			assert.ElementsMatch(t, tc.expectedArchived, res)
		})
	}
}

func TestArchiver_Pipe(t *testing.T) {
	testCases := []struct {
		name             string
		mockSeen         int
		mockArchived     []string
		mockTarget       *domain.Target
		expectedTarget   *domain.Target
		expectedArchived []string
	}{
		{
			"empty",
			0,
			[]string{},
			&domain.Target{BaseURL: "https://www.youtube.com/watch?v=O9FnFxoXBmg&t=1s"},
			&domain.Target{BaseURL: "https://www.youtube.com/watch?v=O9FnFxoXBmg&t=1s"},
			[]string{"https://www.youtube.com/watch?v=O9FnFxoXBmg&t=1s"},
		},
		{
			"notSeen",
			0,
			[]string{"https://fakeSeen.com"},
			&domain.Target{BaseURL: "https://notaSeen.com"},
			&domain.Target{BaseURL: "https://notaSeen.com"},
			[]string{"https://fakeSeen.com", "https://notaSeen.com"},
		},
		{
			"alreadySeen",
			1,
			[]string{"https://yesSeen.com"},
			&domain.Target{BaseURL: "https://yesSeen.com"},
			nil,
			[]string{"https://yesSeen.com"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewArchiver()
			res := make([]string, 0)
			for _, k := range tc.mockArchived {
				a.archive[k] = true
			}

			inChan := make(chan *domain.Target)
			outChan := make(chan *domain.Target)
			wg := sync.WaitGroup{}
			wg.Add(tc.mockSeen)

			go a.Pipe(&wg, inChan, outChan)
			inChan <- tc.mockTarget

			select {
			case res := <-outChan:
				assert.Equal(t, tc.expectedTarget, res)
			case <-time.After(1 * time.Second):
				if tc.mockSeen == 0 {
					t.Errorf("%s timeout", tc.name)
				}
			}

			for k := range a.archive {
				res = append(res, k)
			}
			assert.ElementsMatch(t, tc.expectedArchived, res)
		})
	}
}
