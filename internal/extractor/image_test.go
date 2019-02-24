package extractor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestImage_GetImg(t *testing.T) {
	testCases := []struct {
		name        string
		mockContent []byte
		expected    string
	}{
		{
			"regular",
			[]byte(`<img src="yes.png">`),
			"yes.png",
		},
		{
			"emptyLink",
			[]byte(`<img src="">`),
			"",
		},
		{
			"endTag",
			[]byte(`</img>`),
			"",
		},
		{
			"badType",
			[]byte(`<a href="yes.png">`),
			"",
		},
		{
			"malformed",
			[]byte(`<img href="yes.png">`),
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
			assert.Equal(t, tc.expected, GetImg(tokenizer.Token(), tokenType))
		})
	}
}
