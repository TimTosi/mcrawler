package extractor

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractor_formatLink(t *testing.T) {
	testCases := []struct {
		name               string
		mockURL            string
		mockLink           string
		expected           string
		expectedAssertFunc func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			"regular_fullLinkHTTP",
			"http://www.format.com",
			"http://www.format.com/fullLink",
			"http://www.format.com/fullLink",
			assert.Nil,
		},
		{
			"regular_fullLinkHTTPS",
			"http://www.format.com",
			"https://www.format.com/fullLink",
			"https://www.format.com/fullLink",
			assert.Nil,
		},
		{
			"regular_pathOnly",
			"http://www.format.com",
			"/path-only",
			"http://www.format.com/path-only",
			assert.Nil,
		},
		{
			"regular_queryOnly",
			"http://www.format.com",
			"?arg=ok",
			"http://www.format.com/?arg=ok",
			assert.Nil,
		},
		{
			"regular_fragmentOnly",
			"http://www.format.com",
			"#FragmentOnly",
			"http://www.format.com/#FragmentOnly",
			assert.Nil,
		},
		{
			"regular_pathAndQuery",
			"http://www.format.com",
			"/path-query?ok=toto&test=yes",
			"http://www.format.com/path-query?ok=toto&test=yes",
			assert.Nil,
		},
		{
			"regular_pathAndFragment",
			"http://www.format.com",
			"/path-frag#Fragment",
			"http://www.format.com/path-frag#Fragment",
			assert.Nil,
		},
		{
			"regular_queryAndFragment",
			"http://www.format.com",
			"?arg=ok#Fragment",
			"http://www.format.com/?arg=ok#Fragment",
			assert.Nil,
		},
		{
			"regular_pathAndQueryAndFragment",
			"http://www.format.com",
			"/path-total?arg=ok#Fragment",
			"http://www.format.com/path-total?arg=ok#Fragment",
			assert.Nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := formatLink(tc.mockURL, tc.mockLink)
			assert.Equal(t, tc.expected, res)
			tc.expectedAssertFunc(t, err)
		})
	}
}

func TestExtractor_NewExtractor(t *testing.T) {
	testCases := []struct {
		name           string
		mockCheckFuncs []CheckFunc
	}{
		{
			"regular_noCheckFunc",
			[]CheckFunc{},
		},
		{
			"regular_singleCheckFunc",
			[]CheckFunc{GetImg},
		},
		{
			"regular_multipleCheckFuncs",
			[]CheckFunc{GetImg, GetLinkBasic},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() { NewExtractor(tc.mockCheckFuncs...) })
		})
	}
}

func TestExtractor_ExtractLinks(t *testing.T) {
	testCases := []struct {
		name            string
		mockBaseURL     string
		mockContentPath string
		mockCheckFuncs  []CheckFunc
		expectedLinks   []string
	}{
		{
			"basic_one",
			"http://www.basic-one.com",
			"testdata/one.html",
			[]CheckFunc{GetLinkBasic},
			[]string{"http://www.one.com"},
		},
		{
			"basic_empty",
			"http://www.basic-empty.com",
			"testdata/empty.html",
			[]CheckFunc{GetLinkBasic},
			nil,
		},
		{
			"basic_noLink",
			"http://www.basic-nolink.com",
			"testdata/nolink.html",
			[]CheckFunc{GetLinkBasic},
			nil,
		},
		{
			"basic_noFollow",
			"http://www.basic-nofollow.com",
			"testdata/nofollow.html",
			[]CheckFunc{GetLinkBasic},
			[]string{"http://www.basic-nofollow.com/test"},
		},
		{
			"basic_multiple",
			"https://www.basic-multiple.com",
			"testdata/multiple.html",
			[]CheckFunc{GetLinkBasic},
			[]string{
				"https://www.basic-multiple.com/yes",
				"https://www.basic-multiple.com/test",
				"http://www.ok.com",
			},
		},
		{
			"img_one",
			"https://www.img-one.com",
			"testdata/one.html",
			[]CheckFunc{GetImg},
			[]string{},
		},
		{
			"img_empty",
			"https://www.img-empty.com",
			"testdata/empty.html",
			[]CheckFunc{GetImg},
			[]string{},
		},
		{
			"img_noLink",
			"https://www.img-nolink.com",
			"testdata/nolink.html",
			[]CheckFunc{GetImg},
			[]string{},
		},
		{
			"img_noFollow",
			"https://www.img-nofollow.com",
			"testdata/nofollow.html",
			[]CheckFunc{GetImg},
			[]string{},
		},
		{
			"img_multiple",
			"https://www.img-multiple.com",
			"testdata/multiple.html",
			[]CheckFunc{GetImg},
			[]string{"https://www.img-multiple.com/smiley.gif"},
		},
		{
			"noFollow_one",
			"https://www.nofollow-one.com",
			"testdata/one.html",
			[]CheckFunc{GetLinkNoFollow},
			[]string{"http://www.one.com"},
		},
		{
			"noFollow_empty",
			"https://www.nofollow-empty.com",
			"testdata/empty.html",
			[]CheckFunc{GetLinkNoFollow},
			[]string{},
		},
		{
			"noFollow_noLink",
			"https://www.nofollow-nolink.com",
			"testdata/nolink.html",
			[]CheckFunc{GetLinkNoFollow},
			[]string{},
		},
		{
			"noFollow_noFollow",
			"https://www.nofollow-nofollow.com",
			"testdata/nofollow.html",
			[]CheckFunc{GetLinkNoFollow},
			[]string{},
		},
		{
			"noFollow_multiple",
			"https://www.nofollow-multiple.com",
			"testdata/multiple.html",
			[]CheckFunc{GetLinkNoFollow},
			[]string{
				"https://www.nofollow-multiple.com/yes",
				"http://www.ok.com",
			},
		},
		{
			"multiple_one",
			"https://www.multiple-one.com",
			"testdata/one.html",
			[]CheckFunc{GetLinkBasic, GetImg},
			[]string{"http://www.one.com"},
		},
		{
			"multiple_empty",
			"https://www.multiple-empty.com",
			"testdata/empty.html",
			[]CheckFunc{GetLinkBasic, GetImg},
			[]string{},
		},
		{
			"multiple_noLink",
			"https://www.multiple-nolink.com",
			"testdata/nolink.html",
			[]CheckFunc{GetLinkBasic, GetImg},
			[]string{},
		},
		{
			"multiple_noFollow",
			"https://www.multiple-nofollow.com",
			"testdata/nofollow.html",
			[]CheckFunc{GetLinkBasic, GetImg},
			[]string{"https://www.multiple-nofollow.com/test"},
		},
		{
			"multiple_multiple",
			"https://www.multiple-multiple.com",
			"testdata/multiple.html",
			[]CheckFunc{GetLinkBasic, GetImg},
			[]string{
				"https://www.multiple-multiple.com/yes",
				"https://www.multiple-multiple.com/test",
				"http://www.ok.com",
				"https://www.multiple-multiple.com/smiley.gif",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := ioutil.ReadFile(tc.mockContentPath)
			if err != nil {
				t.Errorf("%s: %v", tc.name, err)
			}

			e := NewExtractor(tc.mockCheckFuncs...)
			res := e.ExtractLinks(tc.mockBaseURL, content)
			assert.ElementsMatch(t, tc.expectedLinks, res)
		})
	}
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)

	os.Exit(m.Run())
}
