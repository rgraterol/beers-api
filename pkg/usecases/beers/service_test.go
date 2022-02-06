package beers_test

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rgraterol/beers-api/cmd/api/initializers"
	"github.com/rgraterol/beers-api/pkg/db"
	"github.com/rgraterol/beers-api/pkg/usecases/beers"
	"github.com/stretchr/testify/assert"
)

func init() {
	initializers.MockDatabaseInitializer()
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

func duplicatedbeerMock() beers.Beer {
	return beers.Beer{
		Name:      "Duplicated",
	}
}

func mockBrokenDB() {
	mockDB, _, _ := sqlmock.New()
	db.Gorm, _ = gorm.Open(mysql.New(mysql.Config{Conn: mockDB, SkipInitializeWithVersion: true}), &gorm.Config{})
}

func clearTestDB() {
	db.Gorm.Exec("DELETE FROM beers")
}
