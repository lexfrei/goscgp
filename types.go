package goscgp

// Condition contains:
// Items on market, Price on market,
// Card condition:
//	1 -- M/NM
// 	2 -- PL
// 	3 -- HP
// 	4 -- Damaged
type Condition struct {
	Condition string
	Count     int
	Price     string
}

type Card struct {
	Name       string
	Set        string
	Foil       bool
	Conditions []Condition
}
