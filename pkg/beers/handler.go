package beers

import (
	"fmt"
	"net/http"
)

func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println(ctx)
}