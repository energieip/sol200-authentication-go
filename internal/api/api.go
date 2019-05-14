package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	pkg "github.com/energieip/common-components-go/pkg/service"
	"github.com/energieip/sol200-authentication-go/internal/database"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/romana/rlog"
)

const (
	APIErrorDeviceNotFound = 1
	APIErrorBodyParsing    = 2
	APIErrorDatabase       = 3
	APIErrorInvalidValue   = 4
	APIErrorUnauthorized   = 5
	APIErrorExpiredToken   = 6

	TokenName = "EiPAccessToken"
)

//APIError Message error code
type APIError struct {
	Code    int    `json:"code"` //errorCode
	Message string `json:"message"`
}

type API struct {
	db              database.Database
	EventsToBackend chan map[string]interface{}
	certificate     string
	keyfile         string
	apiPort         string
	apiPassword     string
	apiIP           string
}

//InitAPI start API connection
func InitAPI(db database.Database, conf pkg.ServiceConfig) *API {
	api := API{
		db:              db,
		EventsToBackend: make(chan map[string]interface{}),
		certificate:     conf.ExternalAPI.CertPath,
		keyfile:         conf.ExternalAPI.KeyPath,
		apiPassword:     conf.ExternalAPI.Password,
		apiPort:         conf.ExternalAPI.Port,
		apiIP:           conf.ExternalAPI.IP,
	}
	go api.swagger()
	return &api
}

func (api *API) setDefaultHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}

func (api *API) sendError(w http.ResponseWriter, errorCode int, message string, httpStatus int) {
	errCode := APIError{
		Code:    errorCode,
		Message: message,
	}

	inrec, _ := json.MarshalIndent(errCode, "", "  ")
	rlog.Error(errCode.Message)
	http.Error(w, string(inrec), httpStatus)
}

type APIInfo struct {
	Versions []string `json:"versions"`
}

type APIFunctions struct {
	Functions []string `json:"functions"`
}

func (api *API) getAPIs(w http.ResponseWriter, req *http.Request) {
	api.setDefaultHeader(w)
	versions := []string{"v1.0"}
	apiInfo := APIInfo{
		Versions: versions,
	}
	inrec, _ := json.MarshalIndent(apiInfo, "", "  ")
	w.Write(inrec)
}

func (api *API) getV1Functions(w http.ResponseWriter, req *http.Request) {
	api.setDefaultHeader(w)
	apiV1 := "/v1.0"
	functions := []string{apiV1 + "/authenticate", apiV1 + "/userInfo", apiV1 + "/userAuthorization"}
	apiInfo := APIFunctions{
		Functions: functions,
	}
	inrec, _ := json.MarshalIndent(apiInfo, "", "  ")
	w.Write(inrec)
}

func (api *API) getFunctions(w http.ResponseWriter, req *http.Request) {
	api.setDefaultHeader(w)
	functions := []string{"/versions"}
	apiInfo := APIFunctions{
		Functions: functions,
	}
	inrec, _ := json.MarshalIndent(apiInfo, "", "  ")
	w.Write(inrec)
}

func (api *API) verification(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenValue := ""
		tokenCookie, err := r.Cookie(TokenName)

		if err != nil || tokenCookie == nil {
			//Check header
			authorizationHeader := r.Header.Get("Authorization")
			if authorizationHeader != "" {
				bearerToken := strings.Split(authorizationHeader, " ")
				if len(bearerToken) > 1 {
					tokenValue = bearerToken[1]
				} else {
					tokenValue = authorizationHeader
				}
			}
		} else {
			tokenValue = tokenCookie.Value
		}
		api.setDefaultHeader(w)

		if tokenValue == "" {
			api.sendError(w, APIErrorUnauthorized, "Unauthorized access", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte(api.apiPassword), nil
		})

		switch err.(type) {
		case nil:
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				api.sendError(w, APIErrorUnauthorized, "Unauthorized access", http.StatusUnauthorized)
				return
			}
			var userClaims Claims
			mapstructure.Decode(claims, &userClaims)
			context.Set(r, "decoded", userClaims)
			next(w, r)

		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)

			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				api.sendError(w, APIErrorExpiredToken, "Expired Token", http.StatusUnauthorized)
				return

			default:
				api.sendError(w, APIErrorBodyParsing, "Error reading request body", http.StatusInternalServerError)
				return
			}

		default:
			api.sendError(w, APIErrorBodyParsing, "Error reading request body", http.StatusInternalServerError)
			return
		}
	})
}

func (api *API) swagger() {
	router := mux.NewRouter()
	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("/media/userdata/www/auth/swaggerui/")))
	router.PathPrefix("/swaggerui/").Handler(sh)

	// API v1.0
	apiV1 := "/v1.0"
	router.HandleFunc(apiV1+"/functions", api.getV1Functions).Methods("GET")

	//setup API
	router.HandleFunc(apiV1+"/user/login", api.createToken).Methods("POST")
	router.HandleFunc(apiV1+"/user", api.verification(api.userInfo)).Methods("GET")
	router.HandleFunc(apiV1+"/userAuthorization", api.verification(api.userAuthorization)).Methods("GET")

	//unversionned API
	router.HandleFunc("/versions", api.getAPIs).Methods("GET")
	router.HandleFunc("/functions", api.getFunctions).Methods("GET")

	log.Fatal(http.ListenAndServeTLS(api.apiIP+":"+api.apiPort, api.certificate, api.keyfile, router))
}
