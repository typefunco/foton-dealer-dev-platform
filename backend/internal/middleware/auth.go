package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/utils/jwt"
)

// AuthMiddleware проверяет JWT токен из cookie или Authorization header
func AuthMiddleware(jwtService *jwt.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Получаем токен из cookie
			cookie, err := c.Cookie("auth_token")
			var token string

			if err == nil && cookie != nil {
				token = cookie.Value
			} else {
				// Если нет cookie, проверяем Authorization header
				authHeader := c.Request().Header.Get("Authorization")
				if authHeader != "" {
					parts := strings.Split(authHeader, " ")
					if len(parts) == 2 && parts[0] == "Bearer" {
						token = parts[1]
					}
				}
			}

			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authentication required",
				})
			}

			// Валидируем токен
			claims, err := jwtService.ValidateJWT(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			// Сохраняем информацию о пользователе в контекст
			c.Set("user", claims)
			c.Set("user_login", claims.Login)
			c.Set("user_is_admin", claims.IsAdmin)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

// AdminMiddleware проверяет, что пользователь является администратором
func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isAdmin, ok := c.Get("user_is_admin").(bool)
			if !ok || !isAdmin {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Admin access required",
				})
			}

			return next(c)
		}
	}
}

// OptionalAuthMiddleware проверяет JWT токен, но не требует его наличия
func OptionalAuthMiddleware(jwtService *jwt.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Получаем токен из cookie
			cookie, err := c.Cookie("auth_token")
			var token string

			if err == nil && cookie != nil {
				token = cookie.Value
			} else {
				// Если нет cookie, проверяем Authorization header
				authHeader := c.Request().Header.Get("Authorization")
				if authHeader != "" {
					parts := strings.Split(authHeader, " ")
					if len(parts) == 2 && parts[0] == "Bearer" {
						token = parts[1]
					}
				}
			}

			if token != "" {
				// Валидируем токен
				claims, err := jwtService.ValidateJWT(token)
				if err == nil {
					// Сохраняем информацию о пользователе в контекст
					c.Set("user", claims)
					c.Set("user_login", claims.Login)
					c.Set("user_is_admin", claims.IsAdmin)
					c.Set("user_role", claims.Role)
				}
			}

			return next(c)
		}
	}
}
