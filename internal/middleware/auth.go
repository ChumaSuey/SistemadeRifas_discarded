package middleware

import (
	"context"
	"net/http"
	"strings"

	"raffle-saas/internal/utils"
)

type contextKey string

const (
	UserIDKey  contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)

func Auth(jwt *utils.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w, utils.NewAPIError(401, "missing authorization header", ""))
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				utils.WriteError(w, utils.NewAPIError(401, "invalid authorization header format", ""))
				return
			}

			claims, err := jwt.ValidateAccessToken(parts[1])
			if err != nil {
				utils.WriteError(w, utils.NewAPIError(401, "invalid or expired token", ""))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) string {
	id, _ := ctx.Value(UserIDKey).(string)
	return id
}

func GetUserRole(ctx context.Context) string {
	role, _ := ctx.Value(UserRoleKey).(string)
	return role
}
