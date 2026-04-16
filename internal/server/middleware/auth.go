package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	userIDContextKey    = "userID"
	roleContextKey      = "role"
)

var (
	ErrNoUserID      = errors.New("userID not found in context")
	ErrInvalidUserID = errors.New("userID is of invalid type")
)

func JWTAuthMiddleware(authSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := GetLoggerFromContext(c)

		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			logger.Warn("empty authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})

			return
		}

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			logger.Warn("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})

			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(authSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), jwt.WithExpirationRequired())

		if err != nil || !token.Valid {
			logger.Warn("invalid jwt token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})

			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Error("failed to parse token claims")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to parse claims"})

			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			logger.Warn("user_id not found in token claims or is not a string")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})

			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.Warn("user_id in token is not a valid UUID", zap.String("user_id", userIDStr))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id format in token"})

			return
		}

		logFields := []zap.Field{zap.String("user_id", userID.String())}

		if role, ok := claims["role"].(string); ok {
			c.Set(roleContextKey, role)
			logFields = append(logFields, zap.String("role", role))
		}

		c.Set(loggerContextKey, logger.With(logFields...))
		c.Set(userIDContextKey, userID)

		c.Next()
	}
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		role, ok := GetRole(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})

			return
		}
		if _, exists := allowed[role]; !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})

			return
		}
		c.Next()
	}
}

func GetRole(c *gin.Context) (string, bool) {
	role, ok := c.Get(roleContextKey)
	if !ok {
		return "", false
	}
	roleStr, ok := role.(string)

	return roleStr, ok
}

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	id, exists := c.Get(userIDContextKey)
	if !exists {
		return uuid.Nil, ErrNoUserID
	}

	userID, ok := id.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrInvalidUserID
	}

	return userID, nil
}
