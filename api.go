package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Response struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SuccessClaim struct {
	Name    string `json:"name"`
	PIN     string `json:"pin"`
	Invoice string `json:"invoice"`
}

// not authenticated, if correct pin is provided call returns the SuccessClaim
func ClaimAddress(w http.ResponseWriter, r *http.Request) {
	params := parseParams(r)
	pin, inv, err := SaveName(params.Name, params, params.Pin)
	if err != nil {
		sendError(w, 400, "could not register name: %s", err.Error())
		return
	}

	response := Response{
		Ok:      true,
		Message: fmt.Sprintf("claimed %v@%v", params.Name, s.Domain),
		Data:    SuccessClaim{params.Name, pin, inv},
	}

	// TODO: middleware for responses that adds this header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	params, err := GetName(name)
	if err != nil {
		sendError(w, 400, err.Error())
		return
	}

	// add pin to response because sometimes not saved in database; after first call to /api/v1/claim
	params.Pin = ComputePIN(name)

	response := Response{
		Ok:      true,
		Message: fmt.Sprintf("%v@%v found", params.Name, s.Domain),
		Data:    params,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := parseParams(r)
	name := mux.Vars(r)["name"]

	// if pin not in json request body get it from header
	if params.Pin == "" {
		// TODO: work with Context()?
		params.Pin = r.Header.Get("X-Pin")
	}

	if _, _, err := SaveName(name, params, params.Pin); err != nil {
		sendError(w, 500, err.Error())
		return
	}

	updatedParams, err := GetName(name)
	if err != nil {
		sendError(w, 500, err.Error())
		return
	}

	// return the updated values or just http.StatusCreated?
	response := Response{
		Ok:      true,
		Message: fmt.Sprintf("updated %v@%v parameters", params.Name, s.Domain),
		Data:    updatedParams,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if err := DeleteName(name); err != nil {
		sendError(w, 500, err.Error())
		return
	}

	response := Response{
		Ok:      true,
		Message: fmt.Sprintf("deleted %v@%v", name, s.Domain),
		Data:    nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// authentication middleware
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// exempt /claim from authentication check;
		if strings.HasPrefix(r.URL.Path, "/api/v1/claim") {
			next.ServeHTTP(w, r)
			return
		}

		name := mux.Vars(r)["name"]
		providedPin := r.Header.Get("X-Pin")

		var err error

		if providedPin == "" {
			err = fmt.Errorf("X-Pin header not provided")
			// pin should always be passed in header but search in json request body anyways
			providedPin = parseParams(r).Pin
		}

		if providedPin != ComputePIN(name) {
			err = fmt.Errorf("wrong pin")
		}

		if err != nil {
			sendError(w, 401, "error fetching user: %s", err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

// helpers
func sendError(w http.ResponseWriter, code int, msg string, args ...interface{}) {
	b, _ := json.Marshal(Response{false, fmt.Sprintf(msg, args...), nil})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func parseParams(r *http.Request) *Params {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var params Params
	json.Unmarshal(reqBody, &params)
	return &params
}
