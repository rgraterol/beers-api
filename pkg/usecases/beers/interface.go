package beers

type Interface interface {
	List() ([]Beer, error)
	Create(b *Beer) (*Beer, error)
	Get(id int) (*Beer, error)
	BoxPrice(id int, boxParams *BeerBoxParameters) (*BeerBox, error)
}
