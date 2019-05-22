package api

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/energieip/sol200-authentication-go/internal/database"
	cmap "github.com/orcaman/concurrent-map"
)

const (
	APIErrorDeviceNotFound = 1
	APIErrorBodyParsing    = 2
	APIErrorDatabase       = 3
	APIErrorInvalidValue   = 4
	APIErrorUnauthorized   = 5
	APIErrorExpiredToken   = 6

	TokenName           = "EiPAccessToken"
	TokenExpirationTime = 86400 // in seconds: 1day
)

//APIError Message error code
type APIError struct {
	Code    int    `json:"code"` //errorCode
	Message string `json:"message"`
}

type API struct {
	db              database.Database
	access          cmap.ConcurrentMap
	EventsToBackend chan map[string]interface{}
	certificate     string
	keyfile         string
	apiPort         string
	apiPassword     string
	apiIP           string
	browsingFolder  string
}

type APIInfo struct {
	Versions []string `json:"versions"`
}

type APIFunctions struct {
	Functions []string `json:"functions"`
}

type JwtToken struct {
	Token     string `json:"accessToken"`
	TokenType string `json:"tokenType"`
	ExpireIn  int    `json:"expireIn"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
