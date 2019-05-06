package api

import (
	"encoding/json"
	"log"
	"net/http"

	pkg "github.com/energieip/common-components-go/pkg/service"
	"github.com/energieip/sol200-authentication-go/internal/core"
	"github.com/energieip/sol200-authentication-go/internal/database"
	"github.com/gorilla/mux"
	"github.com/romana/rlog"
)

type InternalAPI struct {
	db              database.Database
	EventsToBackend chan map[string]interface{}
	certificate     string
	keyfile         string
	apiPort         string
	apiPassword     string
	apiIP           string
}

//InitInternalAPI start API connection
func InitInternalAPI(db database.Database, conf pkg.ServiceConfig) *InternalAPI {
	api := InternalAPI{
		db:              db,
		EventsToBackend: make(chan map[string]interface{}),
		certificate:     conf.Certificate,
		keyfile:         conf.Key,
		apiPassword:     conf.APIPassword,
		apiPort:         "1234",
		apiIP:           "127.0.0.1",
	}
	go api.swagger()
	return &api
}

func (api *InternalAPI) setDefaultHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}

func (api *InternalAPI) sendError(w http.ResponseWriter, errorCode int, message string, httpStatus int) {
	errCode := APIError{
		Code:    errorCode,
		Message: message,
	}

	inrec, _ := json.MarshalIndent(errCode, "", "  ")
	rlog.Error(errCode.Message)
	http.Error(w, string(inrec), httpStatus)
}

func (api *InternalAPI) createInternalUser(w http.ResponseWriter, req *http.Request) {
	api.setDefaultHeader(w)
	var creds core.User
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		api.sendError(w, APIErrorBodyParsing, "Error reading request body", http.StatusInternalServerError)
		return
	}
	event := make(map[string]interface{})
	event[core.CreateUserEvent] = creds
	api.EventsToBackend <- event
	w.Write([]byte("{}"))

}

func (api *InternalAPI) swagger() {
	router := mux.NewRouter()
	router.HandleFunc("/user", api.createInternalUser).Methods("POST")
	log.Fatal(http.ListenAndServeTLS(api.apiIP+":"+api.apiPort, api.certificate, api.keyfile, router))
}
