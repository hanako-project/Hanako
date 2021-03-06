package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/luqmansen/hanako/api/models/postgres"
	"github.com/luqmansen/hanako/api/utils"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//	list of endpoints that doesn't require auth
		noAuth := []string{"/api/v1/user/new", "/api/v1/user/login", "/api/v1/anime/all", "/api/v1/anime/*"}
		requestPath := r.URL.Path //Current request path

		for _, value := range noAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		var response map[string]interface{}
		tokenHeader := r.Header.Get("Authorization") //Grab token from header

		// if token is missing, return 403
		if tokenHeader == "" {
			response = utils.Message(http.StatusForbidden, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		////The token normally comes in format `Bearer {token-body}`,
		// we check if the retrieved token matched this requirement
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = utils.Message(http.StatusForbidden, "Invalid/Malformed Token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return

		}

		tokenPart := splitted[1]
		tk := &postgres.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		//Marformed token, return 403
		if err != nil {
			response = utils.Message(http.StatusForbidden, "Invalid/Malformed Token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}
		//Token is invalid, maybe not signed on this server
		if !token.Valid {
			response = utils.Message(http.StatusForbidden, "Invalid Token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		//Everything good, proceed the request and set the caller to the user retrieved
		//form the parsed token
		_ = fmt.Sprintf("%i", tk.UserId)
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}
