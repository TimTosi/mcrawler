package domain

// Target is a `struct` representing the address of web page to scrape and its
// content.
type Target struct {
	BaseURL string
	Content []byte
	Done    bool
}

// NewTarget returns a new `*domain.Target`.
func NewTarget(baseURL string) *Target {
	return &Target{BaseURL: baseURL, Done: false}
}
