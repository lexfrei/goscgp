package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lexfrei/goscgp"
)

var reStock = regexp.MustCompile(`(?m)(\d+) in stock`)

func ParceCardsInDoc(doc *goquery.Document) ([]goscgp.Card, error) {
	var cards []goscgp.Card
	var breakErr error

	doc.Find("#content > table:nth-child(2) > tbody > tr > td > table > tbody > tr > td").EachWithBreak(
		func(i int, s *goquery.Selection) bool {
			var card goscgp.Card
			card.Name = s.Find("h2").Text()
			card.Set = strings.TrimSuffix(s.Find("div > div.card_desc_details > div:nth-child(1) > div:nth-child(1) > h3 > a").Text(), "\n")

			s.Find("div > div:nth-child(2) > div:nth-child(2) > span").EachWithBreak(func(i int, s *goquery.Selection) bool {
				var con goscgp.Condition

				con.Condition = s.Find("a").Text()[0:2]
				con.Price = s.Find("div:nth-child(4)").Text()
				if reStock.FindStringSubmatch(s.Text()) != nil {
					var err error
					con.Count, err = strconv.Atoi(reStock.FindStringSubmatch(s.Text())[1])
					if err != nil {
						con.Count = 0
						breakErr = err
						return false
					}
				} else {
					con.Count = 0
				}
				card.Conditions = append(card.Conditions, con)
				return true
			})
			if breakErr != nil {
				return false
			}
			cards = append(cards, card)
			return true
		})

	return cards, breakErr
}
