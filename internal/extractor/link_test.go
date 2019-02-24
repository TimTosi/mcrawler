package extractor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestLink_GetLinkBasic(t *testing.T) {
	testCases := []struct {
		name        string
		mockContent []byte
		expected    string
	}{
		{
			"regular",
			[]byte(`<a href="https://www.youtube.com/watch?v=4D2qcbu26gs">`),
			"https://www.youtube.com/watch?v=4D2qcbu26gs",
		},
		{
			"emptyLink",
			[]byte(`<a href="">`),
			"",
		},
		{
			"endTag",
			[]byte(`</a>`),
			"",
		},
		{
			"badType",
			[]byte(`<img src="yes.png">`),
			"",
		},
		{
			"malformed",
			[]byte(`<a src="https://test.com">`),
			"",
		},
		{
			"noFollow",
			[]byte(`<a href="https://still-follow.com" rel="nofollow">`),
			"https://still-follow.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenizer := html.NewTokenizer(bytes.NewReader(tc.mockContent))
			tokenType := tokenizer.Next()
			if tokenType == html.ErrorToken {
				t.Errorf("%s: error token", tc.name)
			}
			assert.Equal(t, tc.expected, GetLinkBasic(tokenizer.Token(), tokenType))
		})
	}
}

func TestLink_GetLinkNoFollow(t *testing.T) {
	testCases := []struct {
		name        string
		mockContent []byte
		expected    string
	}{
		{
			"regular",
			[]byte(`<a href="https://www.youtube.com/watch?v=M9EjE4qm7b8">`),
			"https://www.youtube.com/watch?v=M9EjE4qm7b8",
		},
		{
			"emptyLink",
			[]byte(`<a href="">`),
			"",
		},
		{
			"endTag",
			[]byte(`</a>`),
			"",
		},
		{
			"badType",
			[]byte(`<img src="yes.png">`),
			"",
		},
		{
			"malformed",
			[]byte(`<a src="https://test.com">`),
			"",
		},
		{
			"noFollow",
			[]byte(`<a href="https://no-follow.com" rel="nofollow">`),
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenizer := html.NewTokenizer(bytes.NewReader(tc.mockContent))
			tokenType := tokenizer.Next()
			if tokenType == html.ErrorToken {
				t.Errorf("%s: error token", tc.name)
			}
			assert.Equal(t, tc.expected, GetLinkNoFollow(tokenizer.Token(), tokenType))
		})
	}
}
