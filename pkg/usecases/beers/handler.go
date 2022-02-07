package beers

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rgraterol/beers-api/pkg/responses"
)

const (
	defaultBeerIDParam  = "beerID"
	defaultBeerQuantity = 6
	currencySize        = 3
)
func List(s Interface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		beers, err := s.List()
		if err != nil {
			responses.Error(w, err)
			return
		}
		responses.OK(w, beers)
	}
}

func Create(s Interface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := decodeAndValidateCreateBeerBody(r)
		if err != nil {
			zap.S().Error(err)
			responses.BadRequest(w, err.Error())
			return
		}

		createdB, err := s.Create(b)
		if err == DuplicatedError {
			responses.Duplicated(w, err.Error())
			return
		}
		if err != nil {
			responses.Error(w, err)
			return
		}
		responses.Created(w, createdB)
	}
}

func Get(s Interface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		beerId, err := strconv.Atoi(chi.URLParam(r, defaultBeerIDParam))
		if err != nil {
			zap.S().Error(err)
			responses.BadRequest(w, "invalid " + defaultBeerIDParam)
			return
		}
		beer, err := s.Get(beerId)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			responses.NotFound(w, "beer not found")
			return
		}
		if err != nil {
			responses.Error(w, err)
			return
		}
		responses.OK(w, beer)
	}
}

func BoxPrice(s Interface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		beerId, err := strconv.Atoi(chi.URLParam(r, defaultBeerIDParam))
		if err != nil {
			zap.S().Error(err)
			responses.BadRequest(w, "invalid " + defaultBeerIDParam)
			return
		}
		boxParams, err := decodeBeerBoxPriceParams(r)
		if err != nil {
			zap.S().Error(err)
			responses.BadRequest(w, err.Error())
			return
		}
		beerBox, err := s.BoxPrice(beerId, boxParams)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			responses.NotFound(w, "beer not found")
			return
		}
		if err != nil {
			responses.Error(w, err)
			return
		}
		responses.OK(w, beerBox)
		return
	}
}

func decodeAndValidateCreateBeerBody(r *http.Request) (*Beer, error) {
	var b Beer
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		return nil, err
	}
	if b.Name == "" {
		err = errors.New("name cannot be empty")
		return nil, err
	}
	if b.Price == 0  {
		err = errors.New("price cannot be zero nor empty")
		return nil, err
	}
	if b.Currency == "" || len(b.Currency) != currencySize {
		err = errors.New("currency cannot be empty or different than 3 characters")
		return nil, err
	}
	return &b, err
}

func decodeBeerBoxPriceParams(r *http.Request) (*BeerBoxParameters, error) {
	q, err := strconv.Atoi(r.URL.Query().Get("quantity"))
	if err != nil {
		q = defaultBeerQuantity
	}
	c := r.URL.Query().Get("currency")
	if len(c) != 0 && len(c) != currencySize {
		return nil, errors.New("invalid currency")
	}
	return &BeerBoxParameters{
		Quantity: int64(q),
		Currency: c,
	}, nil
}