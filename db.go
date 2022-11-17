package main

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

const SQL_SETUP_VOTES = `CREATE TABLE IF NOT EXISTS votes (
	author TEXT NOT NULL,
	photo TEXT NOT NULL,
	vote INTEGER,
	UNIQUE(author, photo)
)`

func InitDB() {
	conf, err := pgxpool.ParseConfig(DB_URI)
	PanicIfErr(err)

	db, err := pgxpool.ConnectConfig(context.Background(), conf)
	PanicIfErr(err)

	err = db.Ping(context.Background())
	PanicIfErr(err)

	DB = db

	_, err = DBExec(SQL_SETUP_VOTES)
	PanicIfErr(err)
}

func DBExec(sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return DB.Exec(context.Background(), sql, args...)
}

func DBQuery(sql string, args ...interface{}) (pgx.Rows, error) {
	return DB.Query(context.Background(), sql, args...)
}

func DBQueryRow(sql string, args ...interface{}) pgx.Row {
	return DB.QueryRow(context.Background(), sql, args...)
}
