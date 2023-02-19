package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/varjangn/urlsweetner/db"
)

type ContextKey string

func RequireAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("Authorization")
			if err != nil {
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if float64(time.Now().Unix()) > claims["exp"].(float64) {
					http.Error(w, "Token Expired", http.StatusUnauthorized)
					return
				}
				userEmail := claims["sub"].(string)
				user, err := db.DbRepo.GetUser(userEmail)
				if err != nil {
					http.Error(w, "Token Expired", http.StatusUnauthorized)
					return
				}
				// attach user obj to request
				var key ContextKey = "user"
				ctx := context.Background()
				ctx = context.WithValue(ctx, key, user)
				*r = *r.WithContext(ctx)
				f(w, r)
			}
		}
	}
}
