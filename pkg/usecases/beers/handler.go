package beers

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"

	"github.com/rgraterol/beers-api/pkg/responses"
)

func List(s Interface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		beers := s.List()
		responses.OK(w, beers)
	}
}

func Create(s Interface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := decodeAndValidateBeerBody(r)
		if err != nil {
			zap.S().Error(err)
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

func decodeAndValidateBeerBody(r *http.Request) (Beer, error) {
	var b Beer
	var err error
	if err = json.NewDecoder(r.Body).Decode(&b); err != nil {
		return Beer{}, err
	}
	if b.Name == "" {
		err = errors.New("name cannot be empty")
	}
	return b, err
}