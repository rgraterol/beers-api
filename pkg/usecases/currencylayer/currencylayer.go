package currencylayer

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/pkg/errors"
	"github.com/rgraterol/beers-api/pkg/restclient"
)

const (
	basePath        = "http://api.currencylayer.com"
	freeTrialURL    = "/live"
	accessKey       = "c916fbdbdc0e700ccd61560cafc91fe2"
	DefaultCurrency = "USD"
	cacheTTL        = 30
)
// Raw cache that works, not suited for high concurrency enviroments
var Saved time.Time
var RawCache Response

type CurrencyInterface interface {
	GetCurrency() (*Response, error)
}

type ProductiveLayer struct {}

var Layer CurrencyInterface

//TODO: improve client to a connection pool for enhanced performance
func (l *ProductiveLayer) GetCurrency() (*Response, error) {
	now := time.Now()
	var resp Response
	var err error
	if (Saved == (time.Time{})) || (now.Sub(Saved) > cacheTTL * time.Second) {
		resp, err = executeRequest()
		if err != nil {
			zap.S().Error(err)
			return nil, err
		}
		RawCache = resp
		Saved = time.Now()
	} else {
		resp = RawCache
	}
	return &resp, nil
}

func executeRequest() (Response, error) {
	var resp Response
	res, err := restclient.Get(getURL())
	if err != nil {
		return Response{}, err
	}
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return Response{}, errors.Wrap(err, "cannot get currency layer API")
	}
	return resp, nil
}

func getURL() string {
	return fmt.Sprintf(basePath+freeTrialURL+ "?access_key=%s", accessKey)
}
