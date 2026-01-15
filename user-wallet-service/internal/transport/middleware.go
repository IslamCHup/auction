package transport

import (
    "net/http"
    "strings"

    model "user-service/internal/models"
    "user-service/internal/services"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware(jwt services.JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
            return
        }
        token := strings.TrimPrefix(auth, "Bearer ")
        _, role, uid, err := jwt.ParseToken(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }
        c.Set("user_id", uid)
        c.Set("user_role", string(role))
        c.Next()
    }
}

func RequireRoles(roles ...model.Role) gin.HandlerFunc {
    allowed := map[string]struct{}{}
    for _, r := range roles {
        allowed[string(r)] = struct{}{}
    }
    return func(c *gin.Context) {
        role, _ := c.Get("user_role")
        if roleStr, ok := role.(string); ok {
            if _, ok := allowed[roleStr]; ok {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
    }
}
