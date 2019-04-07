package goscgp

import "fmt"

// Card contains info about card: it's name, set, foiling and conditions
type Card struct {
	Name       string       `json:"Name"`
	Set        string       `json:"Set"`
	Foil       bool         `json:"Foil,omitempty"`
	Conditions []Conditions `json:"Conditions"`
}

// Conditions contains info about card's condition: condition code,
// 													availible units and price
type Conditions struct {
	Condition string `json:"Condition"`
	Count     int    `json:"Count"`
	Price     int    `json:"Price"`
	Discount  int    `json:"Discount,omitempty"`
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
	return str + "\n"
}

func (c Conditions) String() string {
	d := ""
	if c.Discount != 0 {
		d = fmt.Sprintf(" (discount price: $%.2f) ", float64(c.Discount)/100)
	}

	return fmt.Sprintf("%s costs $%.2f%s (%d available)", c.Condition, (float64(c.Price) / 100), d, c.Count)
}
