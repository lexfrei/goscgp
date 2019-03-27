package parser

import (
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/lexfrei/goscgp"
)

var reStock = regexp.MustCompile(`(?m)(\d+) in stock`)
var reNonDigit = regexp.MustCompile(`(?m)\D`)

func ParceCardsOnPage(doc *goquery.Document) ([]goscgp.Card, error) {
	var cards []goscgp.Card

	doc.Find("#content > table:nth-child(2) > tbody > tr > td > table > tbody > tr > td").Each(func(i int, s *goquery.Selection) {
		var card goscgp.Card

		card.Name = s.Find("h2").Text()
		set := strings.TrimSuffix(s.Find("div > div.card_desc_details > div:nth-child(1) > div:nth-child(1) > h3 > a").Text(), "\n")
		if strings.Contains(set, " (Foil)") {
			card.Foil = true
			set = strings.TrimSuffix(set, " (Foil)")
		}
		card.Set = set
		s.Find("div > div:nth-child(2) > div:nth-child(2) > span").Each(func(i int, s *goquery.Selection) {
			var con goscgp.Conditions
			con.Condition = s.Find("a").Text()[0:2]
			p, err := strconv.Atoi(reNonDigit.ReplaceAllString(s.Find("div:nth-child(4)").Text(), ""))
			if err != nil {
				con.Price = 0
			} else {
				con.Price = p
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

func ParceQuery(u url.URL) ([]goscgp.Card, error) {
	var cards []goscgp.Card

	wg := &sync.WaitGroup{}

	jobs := make(chan url.URL, 100)
	results := make(chan []goscgp.Card, 100)

	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go worker(w, jobs, results, wg)
	}

	jobs <- u

	wg.Wait()

	close(results)

	for res := range results {
		cards = append(cards, res...)
	}

	return cards, nil
}

func worker(id int, jobs chan url.URL, results chan<- []goscgp.Card, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		doc, err := goquery.NewDocument(j.String())
		if err != nil {
			os.Exit(1)
		}

		u, exists := doc.Find("#content > table:nth-child(1) > tbody > tr:nth-child(2) > td > div:nth-child(1) > a:contains(\"Next>>\")").Attr("href")

		if exists {
			job, err := url.Parse(u)
			if err != nil {
				os.Exit(51)
			}
			jobs <- *job
		} else {
			close(jobs)
		}

		result, err := ParceCardsOnPage(doc)
		if err != nil {
			os.Exit(1)
		}

		results <- result
	}
}
