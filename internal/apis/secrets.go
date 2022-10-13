package apis

import (
	"app-ez-pwd/internal/secrets"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func RouteUserSecretsApiHandlers(group *echo.Group) {
	group.GET("/categories", ListCategorySecretsGET)
	group.GET("/user-secrets", ListUserSecretsGET)
	group.GET("/user-secrets/:secretId", GetTheUserSecretGET)
	group.POST("/user-secrets", NewUserSecretPOST)
	group.PUT("/user-secrets", UpdateUserSecretsPUT)
	group.DELETE("/user-secrets/:secretId", DeleteUserSecretDELETE)

	group.GET("/user-secrets/backup", GenerateBackupUserSecretsGET)
}

func ListCategorySecretsGET(ctx echo.Context) error {
	rawUserId := ctx.Get("userId")
	userId, _ := rawUserId.(int)

	items, _ := secrets.ListCategorySecretsDB(userId)
	return ctx.JSON(http.StatusOK, items)
}

func ListUserSecretsGET(ctx echo.Context) error {
	rawUserId := ctx.Get("userId")
	userId, _ := rawUserId.(int)

	rawCategoryId := ctx.QueryParam("categoryId")
	categoryId, _ := strconv.ParseInt(rawCategoryId, 10, 32)

	itemsUserSecrets, err := secrets.ListUserSecretDB(userId, int(categoryId))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, itemsUserSecrets)
}

func GetTheUserSecretGET(ctx echo.Context) error {
	rawSecretId := ctx.Param("secretId")

	secretId, err := strconv.ParseInt(rawSecretId, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	rawUserId := ctx.Get("userId")
	userId, _ := rawUserId.(int)

	theSecret, _ := secrets.GetUserSecretByIdDB(userId, int(secretId))
	return ctx.JSON(http.StatusOK, theSecret)
}

func NewUserSecretPOST(ctx echo.Context) error {
	var form secrets.UserSecretForm
	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	if err := form.ValidateFront(); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	rawUserId := ctx.Get("userId")
	userId, _ := rawUserId.(int)

	newSecretId, err := form.Save(userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, map[string]int{"id": newSecretId})
}

func UpdateUserSecretsPUT(ctx echo.Context) error {
	var form secrets.UpdateUserSecretForm
	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	if err := form.ValidateFront(); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	rawUserId := ctx.Get("userId")
	userId, _ := rawUserId.(int)

	err := form.Update(userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, map[string]string{})
}

func DeleteUserSecretDELETE(ctx echo.Context) error {
	rawSecretId := ctx.Param("secretId")

	secretId, err := strconv.ParseInt(rawSecretId, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	rawUserId := ctx.Get("userId")
	userId, _ := rawUserId.(int)

	err = secrets.DeleteUserSecretDB(userId, int(secretId))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, map[string]string{})
}

func GenerateBackupUserSecretsGET(ctx echo.Context) error {
	rawUserId := ctx.Get("userId")
	userId, _ := rawUserId.(int)

	username, byteSecrets := secrets.QueryUserSecretsForExportAsBackup(userId)

	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=the-%s-secrets.zip", username))
	return ctx.Blob(http.StatusOK, "application/zip", byteSecrets)
}
