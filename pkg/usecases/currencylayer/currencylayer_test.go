package currencylayer_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/rgraterol/beers-api/pkg/restclient"
	"github.com/rgraterol/beers-api/pkg/usecases/currencylayer"
	"github.com/stretchr/testify/assert"
)

func init() {
	restclient.Client = &restclient.MockClient{}
}

const jsonMock = `{"source":"USD","quotes": {"USDARS": 105.356594,"USDCLP": 828.503912,"USDEUR":0.873404,"USDUSD":1}}`
const errorMock = `{"success": false,"error": {"code": 101}}`

func init()  {
	currencylayer.Layer = &currencylayer.ProductiveLayer{}
}

func TestGetCurrencyLayerError(t *testing.T) {
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(errorMock)))
	restclient.GetDoFuncMock = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, errors.New("quota exceded")
	}
	resp, err := currencylayer.Layer.GetCurrency()
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestGetCurrencyLayerOk(t *testing.T) {
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonMock)))
	restclient.GetDoFuncMock = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	resp, err := currencylayer.Layer.GetCurrency()
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, float64(1), resp.Quotes["USDUSD"])
}
