package beers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/rgraterol/beers-api/pkg/usecases/beers"
	"github.com/stretchr/testify/assert"
)

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
	ts := httptest.NewServer(http.HandlerFunc(beers.Create(&ServiceMock4XXError{})))
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

func TestList200(t *testing.T) {
	//GIVEN
	ts := httptest.NewServer(http.HandlerFunc(beers.List(&ServiceMockOk{})))
	defer ts.Close()
	//WHEN
	res, _ := http.Get(ts.URL)
	//THEN
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestListError500(t *testing.T) {
	//GIVEN
	ts := httptest.NewServer(http.HandlerFunc(beers.List(&ServiceMockError{})))
	defer ts.Close()
	//WHEN
	res, _ := http.Get(ts.URL)
	//THEN
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGet400(t *testing.T) {
	///GIVEN
	handler := beers.Get(&ServiceMockOk{})
	ts := httptest.NewServer(http.HandlerFunc(handler))
	req := buildRecorderWithContext("", ts.URL)
	w := httptest.NewRecorder()
	//WHEN
	handler(w, req)
	res := w.Result()
	var resp map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&resp)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "invalid beerID", resp["message"])
}

func TestGet200(t *testing.T) {
	///GIVEN
	handler := beers.Get(&ServiceMockOk{})
	ts := httptest.NewServer(http.HandlerFunc(handler))
	req := buildRecorderWithContext("1", ts.URL)
	w := httptest.NewRecorder()
	//WHEN
	handler(w, req)
	res := w.Result()
	var resp map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&resp)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, float64(1), resp["id"])
}

func TestGet404(t *testing.T) {
	///GIVEN
	handler := beers.Get(&ServiceMock4XXError{})
	ts := httptest.NewServer(http.HandlerFunc(handler))
	req := buildRecorderWithContext("1", ts.URL)
	w := httptest.NewRecorder()
	//WHEN
	handler(w, req)
	res := w.Result()
	var resp map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&resp)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t,"beer not found", resp["message"])
}

func buildRecorderWithContext(beerID string, url string) (*http.Request) {
	req := httptest.NewRequest("GET", url, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("beerID", beerID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	return req
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

func (s *ServiceMockOk) List() ([]beers.Beer, error) {
	return nil, nil
}

func (s *ServiceMockOk) Create(b *beers.Beer) (*beers.Beer, error) {
	return &beers.Beer{ID: 1, Name: "test beer"}, nil
}

func (s *ServiceMockOk) Get(id int) (*beers.Beer, error) {
	return &beers.Beer{ID: 1}, nil
}

type ServiceMockError struct {}

func (s *ServiceMockError) List() ([]beers.Beer, error) {
	return nil, errors.New("database connection lost")
}

func (s *ServiceMockError) Create(b *beers.Beer) (*beers.Beer, error) {
	return &beers.Beer{}, errors.New("cannot create new beer")
}

func (s *ServiceMockError) Get(id int) (*beers.Beer, error) {
	return nil, errors.New("cannot get beer")
}

type ServiceMock4XXError struct {}

func (s *ServiceMock4XXError) List() ([]beers.Beer, error) {
	return nil, nil
}

func (s *ServiceMock4XXError) Create(b *beers.Beer) (*beers.Beer, error) {
	return &beers.Beer{}, beers.DuplicatedError
}

func (s *ServiceMock4XXError) Get(id int) (*beers.Beer, error) {
	return nil, gorm.ErrRecordNotFound
}