package middleware

import (
	"net/http"
	"strings"
	"time"

	"example.com/games/environment"
	"github.com/golang-jwt/jwt/v5"
)

type tokenstruct struct {
	Email string
	jwt.RegisteredClaims
}

var jwtKey = []byte(environment.ViperEnvVariable("JWT_TOKEN"))

func CreateToken(email string, w http.ResponseWriter) string {
	expirationTime := time.Now().Add(240 * time.Hour)
	jt := &tokenstruct{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jt)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return ""
	}
	return tokenString

}

func VerifyToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
		tkn, err := jwt.ParseWithClaims(tokenString, &tokenstruct{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
