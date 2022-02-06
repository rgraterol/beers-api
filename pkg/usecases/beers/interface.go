package beers

type Interface interface {
	List() []Beer
	Create(b *Beer) (*Beer, error)
}
