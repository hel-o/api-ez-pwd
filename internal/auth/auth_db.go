package auth

import (
	"app-ez-pwd/internal/logger"
	"app-ez-pwd/internal/storage"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"strings"
)

func GetPasswordHashDB(username string) (int, string, error) {
	username = strings.ToUpper(username)
	query, queryArgs, _ := storage.ApplicationDB.Psql.
		Select("id", "password_hash").
		From("users").
		Where(sq.Eq{
			"UPPER(username)": username,
		}).ToSql()

	cn, tx, _ := storage.ApplicationDB.Begin()
	defer storage.ApplicationDB.Rollback(cn, tx)

	var userId int
	var passwordHash string
	if err := tx.QueryRow(context.Background(), query, queryArgs...).Scan(&userId, &passwordHash); err != nil {
		if err == pgx.ErrNoRows {
			return 0, "", nil
		}
		logger.Logger.Error("err scan", zap.Error(err))
		return 0, "", err
	}

	return userId, passwordHash, nil
}

func ExistsUsernameDB(username string) (bool, error) {
	username = strings.ToUpper(username)
	query, queryArgs, _ := storage.ApplicationDB.Psql.Select("id").
		From("users").
		Where(sq.Eq{
			"UPPER(username)": username,
		}).ToSql()

	cn, tx, _ := storage.ApplicationDB.Begin()
	defer storage.ApplicationDB.Rollback(cn, tx)

	var userId int
	if err := tx.QueryRow(context.Background(), query, queryArgs...).Scan(&userId); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		logger.Logger.Error("err scan", zap.Error(err))
		return false, err
	}

	return userId > 0, nil
}

func SaveNewAccount(username, passwordHash string) error {
	insertUser, insertUserArgs, _ := storage.ApplicationDB.Psql.Insert("users").
		SetMap(map[string]interface{}{
			"username":      username,
			"password_hash": passwordHash,
		}).ToSql() // .Suffix("RETURNING id")

	cn, tx, _ := storage.ApplicationDB.Begin()

	if _, err := tx.Exec(context.Background(), insertUser, insertUserArgs...); err != nil {
		logger.Logger.Error("err insert", zap.Error(err))

		_ = storage.ApplicationDB.Rollback(cn, tx)
		return err
	}

	err := storage.ApplicationDB.Commit(cn, tx)

	return err
}
