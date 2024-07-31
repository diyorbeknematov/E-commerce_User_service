package middleware

import (
	"auth-service/api/token"
	"auth-service/models"
	"log/slog"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func IsAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Cookie'dan tokenni olish
		tokenString, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, models.Errors{
				Error: "Authorization cookie not found",
			})
			ctx.Abort()
			return
		}

		// JWT tokenni tekshirish va tasdiqlash
		claims, err := token.ExtractClaimsAccess(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, models.Errors{
				Error: "Invalid token",
			})
			ctx.Abort()
			return
		}

		// Foydalanuvchi ma'lumotlarini context ga qo'shish
		ctx.Set("claims", claims)

		ctx.Next()
	}
}

func Authorize(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, models.Errors{
				Error: "Unauthorized",
			})
			ctx.Abort()
			return
		}

		userClaims := claims.(*token.Claims)
		ok, err := enforcer.Enforce(userClaims.Role, ctx.Request.URL.Path, ctx.Request.Method)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Errors{
				Error: "Internal server error",
			})
			ctx.Abort()
			return
		}

		if !ok {
			ctx.JSON(http.StatusForbidden, models.Errors{
				Error: "Forbidden",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func LogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Request received",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)

		c.Next()

		logger.Info("Response sent",
			slog.Int("status", c.Writer.Status()),
		)
	}
}