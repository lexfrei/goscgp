package qgen

import "net/url"

type Qurl struct {
	URL *url.URL
}

func NewQurl() *Qurl {
	u, _ := url.Parse("https://www.starcitygames.com/search.php")
	return &Qurl{URL: u}
}

func (q *Qurl) QUrlWithName(name string) {
	qt := q.URL.Query()
	qt.Set("search_query", name)
	q.URL.RawQuery = qt.Encode()
}
