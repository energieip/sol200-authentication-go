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

func (api *API) createToken(w http.ResponseWriter, req *http.Request) {
	api.setDefaultHeader(w, req)
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

	expirationTime := time.Now().Add(TokenExpirationTime * time.Second)
	claims := &Claims{
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

	http.SetCookie(w, &http.Cookie{
		Name:     TokenName,
		Value:    tokenString,
		Expires:  expirationTime,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		MaxAge:   TokenExpirationTime,
	})

	res := JwtToken{
		Token:     tokenString,
		TokenType: "bearer",
		ExpireIn:  TokenExpirationTime,
	}
	api.access.Set(tokenString, *user)
	json.NewEncoder(w).Encode(res)
}

func (api *API) userInfo(w http.ResponseWriter, r *http.Request) {
	decoded := context.Get(r, "decoded")
	var userClaims core.User
	mapstructure.Decode(decoded.(core.User), &userClaims)

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
	var userClaims core.User
	mapstructure.Decode(decoded.(core.User), &userClaims)
	user := database.GetUser(api.db, userClaims.Username)

	if user == nil {
		api.sendError(w, APIErrorDeviceNotFound, "Error information not found", http.StatusInternalServerError)
		return
	}
	permissions := core.UserAuthorization{
		Priviledge:   user.Priviledge,
		AccessGroups: user.AccessGroups,
		Services:     user.Services,
	}
	json.NewEncoder(w).Encode(permissions)
}

func (api *API) logout(w http.ResponseWriter, req *http.Request) {
	decoded := context.Get(req, "token")
	var tokenString string
	mapstructure.Decode(decoded.(string), &tokenString)

	// see https://golang.org/pkg/net/http/#Cookie
	// Setting MaxAge<0 means delete cookie now.

	http.SetCookie(w, &http.Cookie{
		Name:     TokenName,
		MaxAge:   -1,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	})

	api.access.Remove(tokenString)

	w.Write([]byte("{}"))
}
