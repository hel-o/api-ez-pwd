package apis

import (
	"app-ez-pwd/internal/auth"
	"app-ez-pwd/internal/settings"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func DoAuthPOST(ctx echo.Context) error {
	var form auth.DoLoginForm
	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	if err := form.ValidateFront(); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if formErrors := form.Validate(); len(formErrors) > 0 {
		return ctx.JSON(http.StatusBadRequest, formErrors)
	}

	strToken, _ := form.AuthToken()

	secure := !settings.Settings.Debug

	domain := settings.Settings.CookieWebDomain

	userTypeCookie := new(http.Cookie)
	userTypeCookie.Name = "isLogged"
	userTypeCookie.Value = "Y"
	userTypeCookie.Domain = domain
	userTypeCookie.Path = "/"
	userTypeCookie.MaxAge = 0 // Session type
	userTypeCookie.Secure = secure
	userTypeCookie.HttpOnly = false /// readable from javascript frontend
	userTypeCookie.SameSite = http.SameSiteLaxMode

	tokenCookie := new(http.Cookie)
	tokenCookie.Name = "token"
	tokenCookie.Value = strToken
	tokenCookie.Domain = domain
	tokenCookie.Path = "/"
	tokenCookie.MaxAge = 0 // Session type
	tokenCookie.Secure = secure
	tokenCookie.HttpOnly = true
	tokenCookie.SameSite = http.SameSiteLaxMode

	ctx.SetCookie(userTypeCookie)
	ctx.SetCookie(tokenCookie)

	return ctx.JSON(http.StatusOK, map[string]string{})
}

func DoLogoutDELETE(ctx echo.Context) error {
	domain := settings.Settings.CookieWebDomain
	secure := !settings.Settings.Debug
	expire := time.Unix(0, 0)

	userTypeCookie := new(http.Cookie)
	userTypeCookie.Name = "isLogged"
	userTypeCookie.Domain = domain
	userTypeCookie.Expires = expire
	userTypeCookie.Path = "/"
	userTypeCookie.MaxAge = -1
	userTypeCookie.Secure = secure
	userTypeCookie.HttpOnly = false /// readable from javascript frontend
	userTypeCookie.SameSite = http.SameSiteLaxMode

	tokenCookie := new(http.Cookie)
	tokenCookie.Name = "token"
	tokenCookie.Domain = domain
	tokenCookie.Expires = expire
	tokenCookie.Path = "/"
	tokenCookie.MaxAge = -1
	tokenCookie.Secure = secure
	tokenCookie.HttpOnly = true
	tokenCookie.SameSite = http.SameSiteLaxMode

	ctx.SetCookie(userTypeCookie)
	ctx.SetCookie(tokenCookie)

	return ctx.JSON(http.StatusOK, map[string]string{})
}

func CreateNewAccountPOST(ctx echo.Context) error {
	var form auth.NewAccountForm
	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	if err := form.ValidateFront(); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if formErrors := form.Validate(); len(formErrors) > 0 {
		return ctx.JSON(http.StatusBadRequest, formErrors)
	}

	if err := form.Save(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, map[string]string{})
}
