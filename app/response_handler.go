package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dmitriivoitovich/test-assignment-sliide/app/request"
	"github.com/dmitriivoitovich/test-assignment-sliide/app/response"
)

func handleError(w http.ResponseWriter, req *http.Request, err error) {
	status := http.StatusInternalServerError
	if err == request.ErrInvalidParameterValue {
		status = http.StatusBadRequest
	}

	w.WriteHeader(status)
	logRequest(req, status, err)
}

func handleSuccess(w http.ResponseWriter, req *http.Request, resp response.Response) {
	status := http.StatusOK

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logRequest(req, status, err)

		return
	}

	logRequest(req, status, nil)
}

func logRequest(req *http.Request, status int, err error) {
	if err != nil {
		log.Printf("%s %s %d %s", req.Method, req.URL.String(), status, err.Error())

		return
	}

	log.Printf("%s %s %d", req.Method, req.URL.String(), status)
}
