package beers

type Beer struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Brewery  string `json:"brewery"`
	Country  string `json:"chile"`
	Price    string `json:"price"`
	Currency string `json:"currency"`
}

type BeerBox struct {
	ID    int64 `json:"id"`
	Price int64 `json:"price"`
}
