package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func checkAuthInfo(r *http.Request) bool {
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return clientID == os.Getenv("BASIC_AUTH_USER_ID") && clientSecret == os.Getenv("BASIC_AUTH_PASSWORD")
}

func CheckAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthInfo(r) {
			http.Error(w, "Authorization Required", http.StatusUnauthorized)
			fmt.Println("Authorization Required")
		}
		_, err := fmt.Fprintf(w, "Successful Basic Authentication\n")
		if err != nil {
			log.Fatal(err)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
