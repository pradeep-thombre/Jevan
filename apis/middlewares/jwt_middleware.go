package middlewares

import (
	"Jevan/configs"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(configs.AppConfig.JwtSecret),
		TokenLookup: "header:Authorization:Bearer ",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid or missing auth token",
			})
		},
	})
}

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		role, ok := claims["role"].(string)
		fmt.Println("User role:", role)
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Access denied: Admins only",
			})
		}

		return next(c)
	}
}
