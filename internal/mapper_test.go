package crawler

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/timtosi/mcrawler/internal/domain"
)

func TestMapper_NewMapper(t *testing.T) {
	testCases := []struct {
		name               string
		expectedAssertFunc func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{"regular", assert.NotNil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectedAssertFunc(t, NewMapper())
		})
	}
}

func TestMapper_Add(t *testing.T) {
	testCases := []struct {
		name           string
		mockLink       string
		mockSiteMapRaw []string
		expected       []string
	}{
		{
			"empty",
			"https://www.youtube.com/watch?v=DnSXaR5rQcY",
			[]string{},
			[]string{"https://www.youtube.com/watch?v=DnSXaR5rQcY"},
		},
		{
			"notEmpty",
			"https://notaSeen.com",
			[]string{"https://fakeMapper.com"},
			[]string{"https://fakeMapper.com", "https://notaSeen.com"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMapper()
			m.siteMap = tc.mockSiteMapRaw

			assert.NotPanics(t, func() { m.Add(tc.mockLink) })
			assert.ElementsMatch(t, tc.expected, m.siteMap)
		})
	}
}

func TestMapper_Render(t *testing.T) {
	testCases := []struct {
		name           string
		mockSiteMapRaw []string
		expected       string
	}{
		{
			"empty",
			[]string{},
			"testdata/mapper_render_empty.xml",
		},
		{
			"regular_single",
			[]string{"https://fakeMapper.com"},
			"testdata/mapper_render_single.xml",
		},
		{
			"regular_multiple",
			[]string{"https://fakeMapper.com", "https://notaSeen.com"},
			"testdata/mapper_render_multiple.xml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res bytes.Buffer
			stdOut := os.Stdout
			defer func() { os.Stdout = stdOut }()

			r, w, err := os.Pipe()
			if err != nil {
				t.Errorf("%s: %v", tc.name, err)
			}
			os.Stdout = w

			m := NewMapper()
			m.siteMap = tc.mockSiteMapRaw

			assert.NotPanics(t, func() { m.Render() })
			if err := w.Close(); err != nil {
				t.Errorf("%s: %v", tc.name, err)
			}
			io.Copy(&res, r)

			expected, err := ioutil.ReadFile(tc.expected)
			if err != nil {
				t.Errorf("%s: %v", tc.name, err)
			}

			assert.Equal(t, string(expected), res.String())
		})
	}
}

func TestMapper_Pipe(t *testing.T) {
	testCases := []struct {
		name               string
		mockSiteMapRaw     []string
		mockTarget         *domain.Target
		expectedTarget     *domain.Target
		expectedSiteMapRaw []string
	}{
		{
			"empty",
			[]string{},
			&domain.Target{BaseURL: "https://www.youtube.com/watch?v=jooGlIAvDRE"},
			&domain.Target{BaseURL: "https://www.youtube.com/watch?v=jooGlIAvDRE"},
			[]string{"https://www.youtube.com/watch?v=jooGlIAvDRE"},
		},
		{
			"notSeen",
			[]string{"https://fakeMapper.com"},
			&domain.Target{BaseURL: "https://fake.com"},
			&domain.Target{BaseURL: "https://fake.com"},
			[]string{"https://fakeMapper.com", "https://fake.com"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMapper()
			m.siteMap = tc.mockSiteMapRaw

			inChan := make(chan *domain.Target)
			outChan := make(chan *domain.Target)
			wg := sync.WaitGroup{}

			go m.Pipe(&wg, inChan, outChan)
			inChan <- tc.mockTarget

			select {
			case res := <-outChan:
				assert.Equal(t, tc.expectedTarget, res)
				break
			case <-time.After(3 * time.Second):
				t.Errorf("%s timeout", tc.name)
				break
			}

			assert.ElementsMatch(t, tc.expectedSiteMapRaw, m.siteMap)
		})
	}
}
