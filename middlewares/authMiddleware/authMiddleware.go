package authMiddleware

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

const (
	UserEmailKey contextKey = "userEmail"
)

func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/error?message=Unauthorized", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/error?message=Unauthorized", http.StatusSeeOther)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp, ok := claims["exp"].(float64)
			if !ok || time.Now().Unix() > int64(exp) {
				http.Redirect(w, r, "/error?message=Token+expired", http.StatusSeeOther)
				return
			}
		}

		// Extract email from claims
		email := claims["email"].(string)
		ctx := context.WithValue(r.Context(), UserEmailKey, email)
		r = r.WithContext(ctx)

		// Token is valid, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

func TokenCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		// just set it up as empty string for start
		ctx := context.WithValue(r.Context(), UserEmailKey, "")
		r = r.WithContext(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := cookie.Value
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})
		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp, ok := claims["exp"].(float64)
			if !ok || time.Now().Unix() > int64(exp) {
				return
			}
		}

		// Extract email from claims
		email := claims["email"].(string)
		ctx = context.WithValue(r.Context(), UserEmailKey, email)
		r = r.WithContext(ctx)

		// Token is valid, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

