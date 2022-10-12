package storage

import (
	"app-ez-pwd/internal/logger"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type ApplicationDatabase struct {
	ConnString string
	Psql       sq.StatementBuilderType
}

var ApplicationDB ApplicationDatabase

func PrepareApplicationDB(connString string) ApplicationDatabase {
	db := ApplicationDatabase{}
	db.Psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	db.ConnString = connString
	return db
}

func (db ApplicationDatabase) Begin() (*pgx.Conn, pgx.Tx, error) {
	cn, err := pgx.Connect(context.Background(), db.ConnString)
	if err != nil {
		logger.Logger.Error("pgx err connection", zap.Error(err))
		return cn, nil, err
	}
	tx, err := cn.Begin(context.Background())
	return cn, tx, err
}

func (db ApplicationDatabase) Commit(cn *pgx.Conn, tx pgx.Tx) error {
	ctx := context.Background()
	err := tx.Commit(ctx)
	_ = cn.Close(ctx)
	return err
}

func (db ApplicationDatabase) Rollback(cn *pgx.Conn, tx pgx.Tx) error {
	ctx := context.Background()
	err := tx.Rollback(ctx)
	_ = cn.Close(ctx)
	return err
}
