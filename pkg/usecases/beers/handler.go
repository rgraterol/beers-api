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
		b, err := decodeAndValidateBeerBody(r)
		if err != nil {
			responses.BadRequest(w, err.Error())
			return
		}

		createdB, err := s.Create(&b)
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
		strBeerID := chi.URLParam(r, "beerID")
		beerId, err := strconv.Atoi(strBeerID)
		if err != nil {
			zap.S().Error(err)
			responses.BadRequest(w, "invalid beerID")
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

func decodeAndValidateBeerBody(r *http.Request) (Beer, error) {
	var b Beer
	var err error
	if err = json.NewDecoder(r.Body).Decode(&b); err != nil {
		zap.S().Error(err)
		return Beer{}, err
	}
	if b.Name == "" {
		zap.S().Error(err)
		return Beer{}, errors.New("name cannot be empty")
	}
	return b, err
}