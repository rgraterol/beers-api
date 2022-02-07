package beers

import (
	"go.uber.org/zap"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rgraterol/beers-api/pkg/db"
	"github.com/rgraterol/beers-api/pkg/usecases/currencylayer"
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

func (s *Service) BoxPrice(id int, boxParams *BeerBoxParameters) (*BeerBox, error) {
	var box BeerBox
	box.Target = *boxParams
	b, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	box.Beer = *b
	box.Price, err = calculateConvertedPrice(boxParams, b)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	return &box, nil
}

func calculateConvertedPrice(boxParams *BeerBoxParameters, b *Beer) (float64, error) {
	// If two correncies are the same, or doesnt request for a currency conversion
	if boxParams.Currency == "" || boxParams.Currency == b.Currency {
		return float64(boxParams.Quantity) * b.Price, nil
	}
	currencyLayer, err := currencylayer.Layer.GetCurrency()
	if err != nil {
		return 0, errors.Wrap(err, "cannot access currencyLayer API")
	}
	// We get the conversion rate to USD for the requested currency
	usdTarget := currencyLayer.Quotes[currencylayer.DefaultCurrency+ boxParams.Currency]
	if usdTarget == 0 {
		return 0, errors.New("invalid target currency")
	}
	// We get the conversion rate to USD for the beer storage currency
	usdBeer := currencyLayer.Quotes[currencylayer.DefaultCurrency+ b.Currency]
	if usdTarget == 0 {
		return 0, errors.New("invalid beer currency")
	}
	// We get the conversion rate from the beer storage currency to the target one
	conversionRate := usdTarget / usdBeer
	// Finally we multiply the conversion rate with the beer price to get the price in the new currency
	// and we multiply it by the amount of beers in the box
	return b.Price * conversionRate * float64(boxParams.Quantity), nil
}
