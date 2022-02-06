package beers

import (
	"fmt"
	"go.uber.org/zap"
	"strings"

	"github.com/pkg/errors"
	"github.com/rgraterol/beers-api/pkg/db"
)

type Service struct {}

var DuplicatedError = errors.New("the beer already exist in the DB")

func (s *Service) List() ([]Beer, error) {
	var beer Beer
	db.Gorm.First(&beer, "name = ?", "Golden")
	fmt.Println(beer)
	return nil, nil
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
