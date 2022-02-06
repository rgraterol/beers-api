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

func TestCreate200(t *testing.T) {
	// Given
	var s beers.Service
	b := beerMock()
	// When
	_, err := s.Create(&b)
	// Then
	assert.Nil(t, err)
}

func TestCreateDuplicated209(t *testing.T) {
	// Given
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

func TestCreate500(t *testing.T) {
	// Given
	mockDB, _, err := sqlmock.New()
	db.Gorm, err = gorm.Open(mysql.New(mysql.Config{Conn: mockDB, SkipInitializeWithVersion: true}), &gorm.Config{})
	defer initializers.MockDatabaseInitializer()
	var s beers.Service
	b := beerMock()
	// When
	_, err = s.Create(&b)
	// Then
	assert.NotNil(t, err)
	assert.NotEqual(t, beers.DuplicatedError, err)
}

func TestListEmpty200(t *testing.T) {
	// Given
	var s beers.Service
	// When
	bs := s.List()
	// Then
	assert.Equal(t, 0, len(bs))
}

func TestListWithItems200(t *testing.T) {
	// Given
	var s beers.Service
	b := beerMock()
	beerCreate, err := s.Create(&b)
	assert.Nil(t, err)
	// When
	bs := s.List()
	// Then
	assert.Equal(t, bs[0].ID, beerCreate.ID)
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
