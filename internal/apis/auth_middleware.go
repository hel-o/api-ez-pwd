package apis

import (
	"app-ez-pwd/internal/logger"
	"app-ez-pwd/internal/settings"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func VerifyAuthTokenMiddleware(userType string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cookieToken, err := ctx.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					logger.Logger.Warn("request doesn't have the token")
					return ctx.JSON(http.StatusUnauthorized, map[string]string{})
				}
				return err
			}

			tokenUserId, err := evaluateCookieToken(cookieToken.Value)
			if err != nil {
				logger.Logger.Warn("invalid token", zap.Error(err))
				return err
			}

			/*
				// TODO: handle user type, ex: ADMIN
					if tokenUserType != userType {
						return ctx.JSON(http.StatusForbidden, map[string]string{})
					}*/

			ctx.Set("userId", tokenUserId)

			return next(ctx)
		}
	}
}

func evaluateCookieToken(strToken string) (userId int, err error) {
	tokenParser, err := jwt.Parse(strToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "HS256" {
			return nil, errors.New("invalid sign method")
		}
		return []byte(settings.Settings.SecretHex), nil
	})

	if err != nil {
		return 0, err
	}

	if tokenParser.Valid {
		if claims, ok := tokenParser.Claims.(jwt.MapClaims); ok {
			rawUserId, _ := claims["id"].(float64)
			userId = int(rawUserId)

			return userId, nil
		}
	}

	return 0, errors.New("not authenticated")
}
