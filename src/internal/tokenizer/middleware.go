package tokenizer

import (
	"context"
	"ex00/internal/httpserver"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Display struct {
	Token string `json:"tokenizer"`
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var secretJWT = []byte("qwerty") // Create the JWT key used to create the signature

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Date(2022, 6, 9, 12, 0, 0, 0, time.UTC).Unix(),
	}) // Create a new tokenizer object, specifying signing method and the claims you would like it to contain.

	tokenString, err := token.SignedString(secretJWT) // Sign and get the complete encoded tokenizer (string) using the secret
	if err != nil {
		log.Fatal(err)
	}

	tmp := Display{tokenString}

	httpserver.OutputJSON(w, tmp)
}

func Middlewear(next http.Handler) http.Handler {
	secret := []byte("qwerty")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed tokenizer")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "nbf", claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
		}
	})
}
