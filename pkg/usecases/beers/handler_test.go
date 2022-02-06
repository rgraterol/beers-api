package beers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/rgraterol/beers-api/pkg/usecases/beers"
	"github.com/stretchr/testify/assert"
)

func TestList200(t *testing.T) {
	//GIVEN
	ts := httptest.NewServer(http.HandlerFunc(beers.List(&ServiceMockOk{})))
	defer ts.Close()
	//WHEN
	res, _ := http.Get(ts.URL)
	//THEN
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestListError200(t *testing.T) {
	//GIVEN
	ts := httptest.NewServer(http.HandlerFunc(beers.List(&ServiceMockError{})))
	defer ts.Close()
	//WHEN
	res, _ := http.Get(ts.URL)
	//THEN
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestCreateEmptyBody400(t *testing.T) {
	//GIVEN
	ts := httptest.NewServer(http.HandlerFunc(beers.Create(&ServiceMockError{})))
	defer ts.Close()
	//WHEN
	res, _ := http.Post(ts.URL, "application/json", nil)
	var resp map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&resp)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Contains(t, "EOF", resp["message"])
}

func TestCreateEmptyName400(t *testing.T) {
	//GIVEN
	ts := httptest.NewServer(http.HandlerFunc(beers.Create(&ServiceMockError{})))
	defer ts.Close()
	values := map[string]string{}
	body, err := json.Marshal(values)
	assert.Nil(t, err)
	//WHEN
	res, _ := http.Post(ts.URL, "application/json", bytes.NewBuffer(body))
	var resp map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resp)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Contains(t, "name cannot be empty", resp["message"])
}

func TestCreateDuplicated409(t *testing.T) {
	//GIVEN
	ts := httptest.NewServer(http.HandlerFunc(beers.Create(&ServiceMockDuplicatedError{})))
	defer ts.Close()
	//WHEN
	res, _ := http.Post(ts.URL, "application/json", bytes.NewBuffer(buildMockBody()))
	var resp map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&resp)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	assert.Contains(t, beers.DuplicatedError.Error(), resp["message"])
}

func TestCreateError500(t *testing.T) {
	//Given
	ts := httptest.NewServer(http.HandlerFunc(beers.Create(&ServiceMockError{})))
	defer ts.Close()
	//WHEN
	res, _ := http.Post(ts.URL, "application/json", bytes.NewBuffer(buildMockBody()))
	var resp map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&resp)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestCreateOk201(t *testing.T) {
	//GIVEN
	ts := httptest.NewServer(http.HandlerFunc(beers.Create(&ServiceMockOk{})))
	defer ts.Close()
	//WHEN
	res, _ := http.Post(ts.URL, "application/json", bytes.NewBuffer(buildMockBody()))
	var resp map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&resp)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "test beer", resp["name"])
}

func buildMockBody() []byte {
	values := map[string]string{
		"name": "Golden",
	}
	body, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}
	return body
}

type ServiceMockOk struct {}

func (s *ServiceMockOk) List() []beers.Beer {
	return nil
}

func (s *ServiceMockOk) Create(b *beers.Beer) (*beers.Beer, error) {
	return &beers.Beer{ID: 1, Name: "test beer"}, nil
}

type ServiceMockError struct {}

func (s *ServiceMockError) List() []beers.Beer {
	return nil
}

func (s *ServiceMockError) Create(b *beers.Beer) (*beers.Beer, error) {
	return &beers.Beer{}, errors.New("cannot create new beer")
}

type ServiceMockDuplicatedError struct {}

func (s *ServiceMockDuplicatedError) List() []beers.Beer {
	return nil
}

func (s *ServiceMockDuplicatedError) Create(b *beers.Beer) (*beers.Beer, error) {
	return &beers.Beer{}, beers.DuplicatedError
}