package middleware

import (
	"bytes"
	"fmt"
	"food-api/infrastructure/auth"
	"github.com/go-chi/cors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := auth.TokenValid(r)

		if err != nil {
			_ = HTTPError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Basic CORS Support
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"localhost"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		})

		if r.Method == "OPTIONS" {
			_ = HTTPError(w, r, http.StatusNoContent, "")
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//Avoid a large file from loading into memory
//If the file size is greater than 8MB dont allow it to even load into memory and waste our time.
func MaxSizeAllowed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		maxSizeStr := os.Getenv("MAX_SIZE")

		maxSize, err := strconv.ParseInt(fmt.Sprintf(maxSizeStr), 10, 64)
		if err != nil {
			_ = HTTPError(w, r, http.StatusRequestEntityTooLarge, "the maximum size is not an integer")
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		buff, err := ioutil.ReadAll(r.Body)

		if err != nil {
			_ = HTTPError(w, r, http.StatusRequestEntityTooLarge, "too large: upload an image less than 8MB")
			return
		}

		buf := bytes.NewBuffer(buff)
		r.Body = ioutil.NopCloser(buf)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
