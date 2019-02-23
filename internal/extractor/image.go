package extractor

import "golang.org/x/net/html"

// GetImg is an `extractor.CheckFunc` used to retrieve image URLs from a web
// page. It uses `t` as the token to analyse and its `tokenType`. It returns
// the link value or an empty `string` if `t` does not correspond to a link.
func GetImg(t html.Token, tokenType html.TokenType) string {
	if tokenType == html.StartTagToken && t.Data == "img" {
		for _, att := range t.Attr {
			if att.Key == "src" {
				return att.Val
			}
		}
	}
	return ""
}
