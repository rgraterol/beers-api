package beers

import (
	"go.uber.org/zap"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rgraterol/beers-api/pkg/db"
)

type Service struct {}

var DuplicatedError = errors.New("the beer already exist in the DB")

func (s *Service) List() ([]Beer, error) {
	var beers []Beer
	trx := db.Gorm.Find(&beers)
	if trx.Error != nil {
		zap.S().Error("error on list", trx.Error)
		return nil, trx.Error
	}
	return beers, nil
}

func (s *Service) Create(b *Beer) (*Beer, error) {
	trx := db.Gorm.Create(b)
	if trx.Error != nil {
		if strings.Contains(trx.Error.Error(), "Duplicate") || strings.Contains(trx.Error.Error(), "UNIQUE")  {
			zap.S().Error(DuplicatedError, trx.Error)
			return &Beer{}, DuplicatedError
		}
		zap.S().Error("cannot insert beer on DB", trx.Error)
		return &Beer{}, trx.Error
	}
	return b, nil
}

func (s *Service) Get(id int) (*Beer, error) {
	var b Beer
	trx := db.Gorm.First(&b, id)
	if trx.Error != nil {
		zap.S().Error("error getting beer " + strconv.Itoa(id), trx.Error)
		return nil, trx.Error
	}
	return &b, nil
}
