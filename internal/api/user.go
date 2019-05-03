package api

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/energieip/sol200-authentication-go/internal/core"
	"github.com/energieip/sol200-authentication-go/internal/database"
	"github.com/energieip/sol200-authentication-go/internal/tools"
	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
)

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

func (api *API) createToken(w http.ResponseWriter, req *http.Request) {
	api.setDefaultHeader(w)
	var creds Credentials
	err := json.NewDecoder(req.Body).Decode(&creds)

	if err != nil {
		api.sendError(w, APIErrorBodyParsing, "Error reading request body", http.StatusInternalServerError)
		return
	}

	user := database.GetUser(api.db, creds.Username)

	if user == nil || user.Password == nil || !tools.ComparePasswords(*user.Password, creds.Password) {
		api.sendError(w, APIErrorUnauthorized, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(api.apiPassword))
	if err != nil {
		api.sendError(w, APIErrorInvalidValue, "Error during token generation", http.StatusInternalServerError)
		return
	}
	//
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    TokenName,
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })
	res := JwtToken{
		Token:     tokenString,
		TokenType: "bearer",
		ExpireIn:  300,
	}
	json.NewEncoder(w).Encode(res)
}

func (api *API) userInfo(w http.ResponseWriter, r *http.Request) {
	decoded := context.Get(r, "decoded")
	var userClaims Claims
	mapstructure.Decode(decoded.(Claims), &userClaims)

	user := database.GetUser(api.db, userClaims.Username)

	if user == nil {
		api.sendError(w, APIErrorDeviceNotFound, "Error information not found", http.StatusInternalServerError)
		return
	}
	user.Password = nil
	json.NewEncoder(w).Encode(user)
}

func (api *API) userAuthorization(w http.ResponseWriter, r *http.Request) {
	decoded := context.Get(r, "decoded")
	var userClaims Claims
	mapstructure.Decode(decoded.(Claims), &userClaims)
	user := database.GetUser(api.db, userClaims.Username)

	if user == nil {
		api.sendError(w, APIErrorDeviceNotFound, "Error information not found", http.StatusInternalServerError)
		return
	}
	permissions := core.UserAuthorization{
		Priviledges: user.Priviledges,
		AccessGroup: user.AccessGroup,
	}
	json.NewEncoder(w).Encode(permissions)
}
