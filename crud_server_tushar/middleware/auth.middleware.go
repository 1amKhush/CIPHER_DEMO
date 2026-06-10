package middleware

import (
	"context"
	"crud_server/controllers"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// type contextKey string //  creating a new type called contextKey
// const UserContextKey contextKey = "user" // this is done to avoid key collisions, coz it can happen that many have "user"
// //  as their key, but now this is contextKey("user"), this is an IMP CONCEPT

// Not implementing this custom type, and basically this works like when i create a key named "user", there might be
// other req too where i attach claims with the key "user", and thats where key collision would take place, so instead
// of that i just create a new custom type called contextKey, which is nothing but a string, but its use is that suppose
// i have 2 keys named "user", one in auth and one in audit, then the custom type makes sure that "user" of auth != "user" of audit
// bcoz both of them have diff packages, one type is defined in package auth, and the other in is package audit,
// thus it prevents key collision

// Middleware: verify JWT on protected routes
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		decoded_claims := &controllers.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, decoded_claims, getSecret) // 3rd argument is just a func which gives the jwt secret key using which it was encrypted

		// fmt.Printf("The info decoded from the jwt token in the middleware is : %+v\n", decoded_claims)

		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", decoded_claims) // this is creating the new ctx with the claims
		next.ServeHTTP(w, r.WithContext(ctx))                         // now the req has the claims attached to it as ctx
	})
}

func getSecret(token *jwt.Token) (interface{}, error) {
	//ideally you should also be checking the algorithm which was used to encrypt this using token.Method / token.Header["alg"]

	return controllers.JwtSecret, nil
}
