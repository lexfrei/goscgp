package parser

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/lexfrei/goscgp"
)

var reStock = regexp.MustCompile(`(?m)(\d+) in stock`)
var reNonDigit = regexp.MustCompile(`(?m)\D`)

func parceCardsOnPage(doc *goquery.Document) ([]goscgp.Card, error) {
	var cards []goscgp.Card

	doc.Find("#content > table:nth-child(2) > tbody > tr > td > table > tbody > tr > td").Each(
		func(i int, s *goquery.Selection) {
			var card goscgp.Card

			card.Name = s.Find("h2").Text()

			set := strings.TrimSuffix(s.Find("div > div.card_desc_details > div:nth-child(1) > div:nth-child(1) > h3 > a").Text(), "\n")
			if strings.Contains(set, " (Foil)") {
				card.Foil = true
				set = strings.TrimSuffix(set, " (Foil)")
			}
			card.Set = set

			s.Find("div > div:nth-child(2) > div:nth-child(2) > span").Each(
				func(i int, s *goquery.Selection) {
					var con goscgp.Conditions
					con.Condition = s.Find("a").Text()[0:2]

					_, e := s.Find("div:nth-child(4) > span").Attr("style")

					if e {
						p, err := strconv.Atoi(reNonDigit.ReplaceAllString(
							s.Find("div:nth-child(4) > span:nth-child(1)").Text(), ""))
						if err != nil {
							con.Price = 0
						} else {
							con.Price = p
						}

						d, err := strconv.Atoi(reNonDigit.ReplaceAllString(
							s.Find("div:nth-child(4) > span:nth-child(3)").Text(), ""))
						if err != nil {
							con.Discount = 0
						} else {
							con.Discount = d
						}

					} else {
						p, err := strconv.Atoi(reNonDigit.ReplaceAllString(
							s.Find("div:nth-child(4)").Text(), ""))
						if err != nil {
							con.Price = 0
						} else {
							con.Price = p
						}
					}

					if reStock.FindStringSubmatch(s.Text()) != nil {
						var err error
						con.Count, err = strconv.Atoi(reStock.FindStringSubmatch(s.Text())[1])
						if err != nil {
							// FIXME: Should return error
							con.Count = 0
						}
					} else {
						con.Count = 0
					}
					card.Conditions = append(card.Conditions, con)
				})
			cards = append(cards, card)
		})

	return cards, nil
}

func DoRequest(u url.URL, c *http.Client) ([]goscgp.Card, error) {
	var cards []goscgp.Card

	wg := &sync.WaitGroup{}

	jobs := make(chan url.URL, 100)
	results := make(chan []goscgp.Card, 100)

	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go worker(jobs, results, wg, c)
	}

	jobs <- u

	wg.Wait()

	close(results)

	for res := range results {
		cards = append(cards, res...)
	}

	return cards, nil
}

func worker(jobs chan url.URL, results chan<- []goscgp.Card, wg *sync.WaitGroup, c *http.Client) {
	defer wg.Done()
	for j := range jobs {
		doc, err := getDocumentFromURL(c, j)
		if err != nil {
			continue
		}

		u, exists := doc.Find("#content > table:nth-child(1) > tbody > tr:nth-child(2) > td > div:nth-child(1) > a:contains(\"Next>>\")").Attr("href")

		if exists {
			job, err := url.Parse(u)
			if err != nil {
				continue
			}
			jobs <- *job
		} else {
			close(jobs)
		}

		result, err := parceCardsOnPage(doc)
		if err != nil {
			continue
		}

		results <- result
	}
}

func getDocumentFromURL(c *http.Client, url url.URL) (*goquery.Document, error) {
	res, e := c.Get(url.String())
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()

	if res == nil {
		return nil, errors.New("response is nil")
	}

	return goquery.NewDocumentFromReader(res.Body)
}
