package initializers

import (
	"net/http"

	"github.com/rgraterol/beers-api/pkg/restclient"
	"github.com/rgraterol/beers-api/pkg/usecases/currencylayer"
)

func RestClientsInitializer() {
	restclient.Client = &http.Client{}
	currencylayer.Layer = &currencylayer.ProductiveLayer{}
}
