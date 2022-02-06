package beers_test

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/rgraterol/beers-api/cmd/api/initializers"
	"github.com/rgraterol/beers-api/pkg/db"
	"github.com/rgraterol/beers-api/pkg/usecases/beers"
	"github.com/rgraterol/beers-api/pkg/usecases/currencylayer"
	"github.com/stretchr/testify/assert"
)

func init() {
	initTestDB()
}

func TestCreateOk(t *testing.T) {
	// Given
	clearTestDB()
	var s beers.Service
	b := beerMock()
	// When
	_, err := s.Create(&b)
	// Then
	assert.Nil(t, err)
}

func TestCreateDuplicated(t *testing.T) {
	// Given
	clearTestDB()
	var s beers.Service
	b := duplicatedbeerMock()
	// When
	_, err := s.Create(&b)
	assert.Nil(t, err)
	_, err = s.Create(&b)
	// Then
	assert.NotNil(t, err)
	assert.Equal(t, beers.DuplicatedError, err)
}

func TestCreateError(t *testing.T) {
	// Given
	mockBrokenDB()
	defer initializers.MockDatabaseInitializer()
	var s beers.Service
	b := beerMock()
	// When
	_, err := s.Create(&b)
	// Then
	assert.NotNil(t, err)
	assert.NotEqual(t, beers.DuplicatedError, err)
}

func TestListError(t *testing.T) {
	// Given
	mockBrokenDB()
	defer initializers.MockDatabaseInitializer()
	var s beers.Service
	// When
	bs, err := s.List()
	// Then
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(bs))
}

func TestListEmpty(t *testing.T) {
	// Given
	clearTestDB()
	var s beers.Service
	// When
	bs, err := s.List()
	// Then
	assert.Nil(t, err)
	assert.Equal(t, 0, len(bs))
}

func TestListWithItems(t *testing.T) {
	// Given
	clearTestDB()
	var s beers.Service
	b := beerMock()
	beerCreate, err := s.Create(&b)
	assert.Nil(t, err)
	// When
	bs, err := s.List()
	// Then
	assert.Nil(t, err)
	assert.Equal(t, bs[0].ID, beerCreate.ID)
}

func TestGetNotFound(t *testing.T) {
	// Given
	clearTestDB()
	var s beers.Service
	// When
	b, err := s.Get(1)
	// Then
	assert.NotNil(t, err)
	assert.Nil(t,b)
}

func TestGetOK(t *testing.T) {
	// Given
	clearTestDB()
	var s beers.Service
	b := beerMock()
	// When
	_, err := s.Create(&b)
	fetchedB, err := s.Get(1)
	// Then
	assert.Nil(t, err)
	assert.NotNil(t,fetchedB)
	assert.Equal(t, b.ID, fetchedB.ID)
}

func TestBoxPriceNotFoundError(t *testing.T) {
	// Given
	clearTestDB()
	var s beers.Service
	// When
	b, err := s.BoxPrice(1, &beers.BeerBoxParameters{})
	// Then
	assert.NotNil(t, err)
	assert.Nil(t,b)
}

func TestBoxPriceClientLayerError(t *testing.T) {
	// Given
	clearTestDB()
	var s beers.Service
	b := specificPriceBeerMock()
	_, err := s.Create(&b)
	assert.Nil(t, err)
	currencylayer.Layer = &mockLayerError{}
	// When
	p, err := s.BoxPrice(2, &beers.BeerBoxParameters{
		Currency: "NYC",
	})
	// Then
	assert.NotNil(t, err)
	assert.Nil(t,p)
	assert.Contains(t, err.Error(), "cannot access currencyLayer API")
}

func TestBoxPriceInvalidCurrencyError(t *testing.T) {
	// Given
	clearTestDB()
	var s beers.Service
	b := specificPriceBeerMock()
	_, err := s.Create(&b)
	assert.Nil(t, err)
	currencylayer.Layer = &mockLayerOk{}
	// When
	p, err := s.BoxPrice(2, &beers.BeerBoxParameters{
		Currency: "NYC",
	})
	// Then
	assert.NotNil(t, err)
	assert.Nil(t,p)
	assert.Contains(t, err.Error(), "invalid target currency")
}

func TestBoxPriceOkConversion(t *testing.T) {
	// Given
	clearTestDB()
	var s beers.Service
	b := specificPriceBeerMock()
	_, err := s.Create(&b)
	assert.Nil(t, err)
	currencylayer.Layer = &mockLayerOk{}
	// When
	p, err := s.BoxPrice(2, &beers.BeerBoxParameters{
		Quantity: 12,
		Currency: "ARS",
	})
	// Then
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, float64(2288.9676977167974), p.Price)
	assert.Equal(t, b.ID, p.Beer.ID)
	assert.Equal(t, "ARS", p.Target.Currency)
}


func beerMock() beers.Beer {
	return beers.Beer{
		ID:        1,
		Name:      "GoldenMock",
		Brewery:   "MockingBrewery",
		Country:   "ChileMock",
		Price:     2.5,
		Currency:  "MCK",
	}
}

func specificPriceBeerMock() beers.Beer {
	return beers.Beer{
		ID:        2,
		Name:      "Calafate",
		Brewery:   "Austral",
		Country:   "ChileMock",
		Price:     1500,
		Currency:  "CLP",
	}
}

func duplicatedbeerMock() beers.Beer {
	return beers.Beer{
		Name:      "Duplicated",
	}
}

func mockBrokenDB() {
	mockDB, _, _ := sqlmock.New()
	db.Gorm, _ = gorm.Open(mysql.New(mysql.Config{Conn: mockDB, SkipInitializeWithVersion: true}), &gorm.Config{})
}

func initTestDB() {
	var err error
	db.Gorm, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(errors.Wrap(err, "failed to connect gorm with mock DB"))
	}
	db.Gorm.AutoMigrate(&beers.Beer{})
}

func clearTestDB() {
	db.Gorm.Exec("DELETE FROM beers")
}

type mockLayerOk struct{}

func (l *mockLayerOk) GetCurrency() (*currencylayer.Response, error) {
	return &currencylayer.Response{
		Source: "USD",
		Quotes: map[string]float64{
			"USDCLP": float64(828.503912),
			"USDARS": float64(105.356594),
			"USDEUR": float64(0.873404),
			"USDUSD": float64(1),
		},
	}, nil
}

type mockLayerError struct{}

func (l *mockLayerError) GetCurrency() (*currencylayer.Response, error) {
	return nil, errors.New("error with layer")
}
