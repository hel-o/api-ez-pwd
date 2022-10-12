package secrets

import (
	"app-ez-pwd/internal/logger"
	"app-ez-pwd/internal/storage"
	"context"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type ListCategoryModel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func ListCategorySecretsDB(userId int) ([]ListCategoryModel, error) {
	itemsCategory := make([]ListCategoryModel, 0)

	selectQry, selectQryArgs, _ := storage.ApplicationDB.Psql.
		Select("id", "name").
		From("secret_categories").Where(sq.Eq{
		"user_id": userId,
	}).OrderBy("name DESC").ToSql()

	cn, tx, _ := storage.ApplicationDB.Begin()
	defer storage.ApplicationDB.Rollback(cn, tx)

	rows, err := tx.Query(context.Background(), selectQry, selectQryArgs...)
	if err != nil {
		logger.Logger.Error("err query", zap.Error(err))
		return itemsCategory, err
	}

	defer rows.Close()
	for rows.Next() {
		var item ListCategoryModel
		err = rows.Scan(&item.Id, &item.Name)
		if err != nil {
			logger.Logger.Error("err scan", zap.Error(err))
			return itemsCategory, err
		}

		itemsCategory = append(itemsCategory, item)
	}

	return itemsCategory, err
}

type ListUserSecretModel struct {
	Id                int             `json:"id"`
	Description       string          `json:"description"`
	Username          string          `json:"username"`
	PasswordEncrypted json.RawMessage `json:"passwordEncrypted"`
	SafeNoteEncrypted json.RawMessage `json:"safeNoteEncrypted"`
	URLSite           string          `json:"URLSite"`
}

func ListUserSecretDB(userId, categoryId int) ([]ListUserSecretModel, error) {
	itemsUserSecrets := make([]ListUserSecretModel, 0)

	whereFilters := sq.Eq{
		"user_id": userId,
	}

	if categoryId > 0 {
		whereFilters["category_id"] = categoryId
	}

	selectQry, selectQryArgs, _ := storage.ApplicationDB.Psql.Select(
		"id",
		"description",
		"username",
		"password_json",
		"safe_note_json",
		"url_site",
	).From("user_secrets").Where(whereFilters).OrderBy("id DESC").ToSql()

	cn, tx, _ := storage.ApplicationDB.Begin()
	defer storage.ApplicationDB.Rollback(cn, tx)

	rows, err := tx.Query(context.Background(), selectQry, selectQryArgs...)
	if err != nil {
		logger.Logger.Error("err select qry", zap.Error(err))
		return itemsUserSecrets, err
	}

	defer rows.Close()
	for rows.Next() {
		var userSecret ListUserSecretModel
		err = rows.Scan(
			&userSecret.Id,
			&userSecret.Description,
			&userSecret.Username,
			&userSecret.PasswordEncrypted,
			&userSecret.SafeNoteEncrypted,
			&userSecret.URLSite,
		)
		if err != nil {
			logger.Logger.Error("err scan item", zap.Error(err))
			return itemsUserSecrets, err
		}

		itemsUserSecrets = append(itemsUserSecrets, userSecret)
	}

	return itemsUserSecrets, err
}

type UserSecretModel struct {
	Id                int             `json:"id"`
	CategoryId        int             `json:"categoryId"`
	Description       string          `json:"description"`
	Username          string          `json:"username"`
	PasswordEncrypted json.RawMessage `json:"passwordEncrypted"`
	SafeNoteEncrypted json.RawMessage `json:"safeNoteEncrypted"`
	URLSite           string          `json:"urlSite"`
}

func GetUserSecretByIdDB(userId, secretId int) (userSecret UserSecretModel, err error) {
	selectSecret, selectSecretArgs, _ := storage.ApplicationDB.Psql.Select(
		"id",
		"description",
		"username",
		"password_json",
		"safe_note_json",
		"url_site",
		"category_id",
	).From("user_secrets").Where(sq.Eq{
		"user_id": userId,
		"id":      secretId,
	}).ToSql()

	cn, tx, _ := storage.ApplicationDB.Begin()
	defer storage.ApplicationDB.Rollback(cn, tx)

	err = tx.QueryRow(context.Background(), selectSecret, selectSecretArgs...).Scan(
		&userSecret.Id,
		&userSecret.Description,
		&userSecret.Username,
		&userSecret.PasswordEncrypted,
		&userSecret.SafeNoteEncrypted,
		&userSecret.URLSite,
		&userSecret.CategoryId,
	)
	if err != nil {
		logger.Logger.Error("err query", zap.Error(err))
		return userSecret, err
	}
	return userSecret, err
}

type NewUserSecretModel struct {
	UserId            int
	CategoryId        int
	NewCategoryName   string
	Description       string
	Username          string
	PasswordEncrypted []byte
	SafeNoteEncrypted []byte
	URLSite           string
}

func SaveNewUserSecretDB(newUserSecret NewUserSecretModel) (int, error) {
	cn, tx, _ := storage.ApplicationDB.Begin()

	if newUserSecret.CategoryId == 0 {
		insertNewCategory, insertNewCategoryArgs, _ := storage.ApplicationDB.Psql.Insert("secret_categories").
			SetMap(map[string]interface{}{
				"name":    newUserSecret.NewCategoryName,
				"user_id": newUserSecret.UserId,
			}).Suffix("RETURNING id").ToSql()

		err := tx.QueryRow(context.Background(), insertNewCategory, insertNewCategoryArgs...).Scan(&newUserSecret.CategoryId)
		if err != nil {
			logger.Logger.Error("err insert new category", zap.Error(err))

			_ = storage.ApplicationDB.Rollback(cn, tx)

			return 0, err
		}

	}

	insertSecret, insertSecretArgs, _ := storage.ApplicationDB.Psql.Insert("user_secrets").
		SetMap(map[string]interface{}{
			"description":    newUserSecret.Description,
			"username":       newUserSecret.Username,
			"password_json":  newUserSecret.PasswordEncrypted,
			"safe_note_json": newUserSecret.SafeNoteEncrypted,
			"url_site":       newUserSecret.URLSite,
			"category_id":    newUserSecret.CategoryId,
			"user_id":        newUserSecret.UserId,
		}).Suffix("RETURNING id").ToSql()

	var newSecretId int
	err := tx.QueryRow(context.Background(), insertSecret, insertSecretArgs...).Scan(&newSecretId)
	if err != nil {
		logger.Logger.Error("err insert user secret", zap.Error(err))

		_ = storage.ApplicationDB.Rollback(cn, tx)

	} else {
		_ = storage.ApplicationDB.Commit(cn, tx)
	}

	return newSecretId, err
}

type UpdateUserSecretModel struct {
	UserId            int
	SecretId          int
	CategoryId        int
	NewCategoryName   string
	Description       string
	Username          string
	PasswordEncrypted []byte
	SafeNoteEncrypted []byte
	URLSite           string
}

func UpdateUserSecretDB(userSecretModel UpdateUserSecretModel) error {
	cn, tx, _ := storage.ApplicationDB.Begin()

	if userSecretModel.CategoryId == 0 {
		insertNewCategory, insertNewCategoryArgs, _ := storage.ApplicationDB.Psql.Insert("secret_categories").
			SetMap(map[string]interface{}{
				"name":    userSecretModel.NewCategoryName,
				"user_id": userSecretModel.UserId,
			}).Suffix("RETURNING id").ToSql()

		err := tx.QueryRow(context.Background(), insertNewCategory, insertNewCategoryArgs...).Scan(&userSecretModel.CategoryId)
		if err != nil {
			logger.Logger.Error("err insert new category", zap.Error(err))

			_ = storage.ApplicationDB.Rollback(cn, tx)

			return err
		}
	}

	updateUserSecret, updateUserSecretArgs, _ := storage.ApplicationDB.Psql.Update("user_secrets").
		SetMap(map[string]interface{}{
			"description":    userSecretModel.Description,
			"username":       userSecretModel.Username,
			"password_json":  userSecretModel.PasswordEncrypted,
			"safe_note_json": userSecretModel.SafeNoteEncrypted,
			"url_site":       userSecretModel.URLSite,
			"category_id":    userSecretModel.CategoryId,
		}).Where(sq.Eq{
		"id":      userSecretModel.SecretId,
		"user_id": userSecretModel.UserId,
	}).ToSql()

	_, err := tx.Exec(context.Background(), updateUserSecret, updateUserSecretArgs...)
	if err != nil {
		logger.Logger.Error("err updating user secret", zap.Error(err))

		_ = storage.ApplicationDB.Rollback(cn, tx)
	} else {
		_ = storage.ApplicationDB.Commit(cn, tx)
	}

	return err
}

func DeleteUserSecretDB(userId, secretId int) error {
	deleteQry, deleteArgs, _ := storage.ApplicationDB.Psql.Delete("user_secrets").Where(sq.Eq{
		"user_id": userId,
		"id":      secretId,
	}).ToSql()

	cn, tx, _ := storage.ApplicationDB.Begin()

	_, err := tx.Exec(context.Background(), deleteQry, deleteArgs...)
	if err != nil {

		logger.Logger.Error("err deleting secret", zap.Error(err))
		_ = storage.ApplicationDB.Rollback(cn, tx)

	} else {
		_ = storage.ApplicationDB.Commit(cn, tx)
	}

	return err
}
