package goscgp

import "fmt"

type Card struct {
	Name       string       `json:"Name"`
	Set        string       `json:"Set"`
	Foil       bool         `json:"Foil"`
	Conditions []Conditions `json:"Conditions"`
}

// Conditions contains:
// Items on market, Price on market,
// Card condition:
//	1 -- M/NM
// 	2 -- PL
// 	3 -- HP
// 	4 -- Damaged
type Conditions struct {
	Condition string `json:"Condition"`
	Count     int    `json:"Count"`
	Price     int    `json:"Price"`
}

func (c Card) String() string {
	str := "*" + c.Name + "*"
	if c.Foil {
		str = str + " (foil)"
	}
	str = str + " " + c.Set
	for _, v := range c.Conditions {
		str = str + fmt.Sprintf("\t\n%s", v.String())
	}
	return str
}

func (c Conditions) String() string {
	return fmt.Sprintf("%s costs *%.2f (%d availeble)", c.Condition, (float64(c.Price) / 100), c.Count)
}
