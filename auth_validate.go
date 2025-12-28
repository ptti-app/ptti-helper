package ptti

import (
	"context"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

type payload struct {
	Token string `json:"token"`
}

type validateTokenRes struct {
	UserID string `json:"user_id"`
}

type ContextValue struct {
	User string `json:"user"`
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "auth header missing", http.StatusUnauthorized)
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)
		svcURL := GetEnv("AUTH_SERVICE_URL")
		endpoint := GetEnv("AUTH_VALIDATE_ENDPOINT")
		validateTokenURL := svcURL + endpoint

		payload := payload{
			Token: token,
		}

		var response validateTokenRes

		err := PostJSON(validateTokenURL, payload, &response)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if response.UserID == "" {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		userObjID, err := bson.ObjectIDFromHex(response.UserID)
		if err != nil {
			http.Error(w, "invalid objectid", http.StatusUnauthorized)
			return
		}

		userCtx := ContextValue{
			User: userObjID.Hex(),
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
