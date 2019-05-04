package xnethttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amaya382/xerrorz"
)

func SetHTTPErrJSON(w http.ResponseWriter, errType xerrorz.ErrType, innerErrs ...*xerrorz.InnerErr) {
	httpErr := xerrorz.NewHTTPErr(errType, innerErrs...)
	bJSON, err := json.Marshal(httpErr)
	if err != nil {
		panic("Failed to generate an error response")
	}

	// Write
	w.WriteHeader(httpErr.ErrDoc.Code)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(bJSON))
}
