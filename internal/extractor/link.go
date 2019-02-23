package extractor

import "golang.org/x/net/html"

// GetLinkNoFollow is an `extractor.CheckFunc` used to retrieve link URLs from
// a web page. It uses `t` as the token to analyse and its `tokenType`.
// It returns the link value or an empty `string` if `t` does not correspond
// to a link.
//
// NOTE: This function respect the `nofollow` meta tag.
func GetLinkNoFollow(t html.Token, tokenType html.TokenType) string {
	link := ""

	if tokenType == html.StartTagToken && t.Data == "a" {
		for _, att := range t.Attr {
			if att.Key == "href" {
				link = att.Val
			} else if att.Key == "rel" && att.Val == "nofollow" {
				return ""
			}
		}
	}
	return link
}

// GetLinkBasic is an `extractor.CheckFunc` used to retrieve link URLs from a
// web page. It uses `t` as the token to analyse and its `tokenType`.
// It returns the link value or an empty `string` if `t` does not correspond
// to a link.
//
// NOTE: This function ignores the `nofollow` meta tag.
func GetLinkBasic(t html.Token, tokenType html.TokenType) string {
	if tokenType == html.StartTagToken && t.Data == "a" {
		for _, att := range t.Attr {
			if att.Key == "href" {
				return att.Val
			}
		}
	}
	return ""
}
