package resp

import (
	"encoding/json"
	"net/http"
)

type resp struct {
	Data  *interface{} `json:"data"`
	Error string       `json:"error,omitempty"`
}

func WriteResp(w http.ResponseWriter, data interface{}, err error) {
	var r resp
	if err == nil {
		r.Data = &data
	} else {
		r.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(r)
}
