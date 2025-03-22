package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool
var Q *Queries

func Init(connString string) error {
	migrator, err := migrate.New("file://migrations", connString+"?sslmode=disable")
	if err != nil {
		return err
	}

	fmt.Println("Migrating DB")
	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	conn, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		return err
	}

	Q = New(conn)

	return nil
}

func NewTx(ctx context.Context) (pgx.Tx, error) {
	return conn.Begin(ctx)
}
