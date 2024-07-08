package storage

import (
	"fmt"
	"log"

	"github.com/crt379/svc-collector-grpc/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	WriteDB *sqlx.DB
	ReadDB  *sqlx.DB
)

func init() {
	var err error

	WriteDB, err = NewPgConnect(
		config.AppConfig.PgSql.Write.Host,
		config.AppConfig.PgSql.Write.Port,
		config.AppConfig.PgSql.Write.User,
		config.AppConfig.PgSql.Write.Password,
		config.AppConfig.PgSql.Write.Dbname,
	)
	if err != nil {
		log.Panicf(err.Error())
	}

	WriteDB.SetMaxOpenConns(20)
	WriteDB.SetMaxIdleConns(10)

	if config.AppConfig.PgSql.Read.Port == "" {
		ReadDB = WriteDB
	} else {
		ReadDB, err = NewPgConnect(
			config.AppConfig.PgSql.Read.Host,
			config.AppConfig.PgSql.Read.Port,
			config.AppConfig.PgSql.Read.User,
			config.AppConfig.PgSql.Read.Password,
			config.AppConfig.PgSql.Read.Dbname,
		)
		if err != nil {
			log.Panicf(err.Error())
		}

		ReadDB.SetMaxOpenConns(20)
		ReadDB.SetMaxIdleConns(10)
	}
}

func NewPgConnect(host, port, user, password, database string) (*sqlx.DB, error) {
	s := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, database)
	db, err := sqlx.Connect("pgx", s)
	if err != nil {
		log.Printf("failed to connect to postgres, host:%s port:%s user:%s err: %s\n", host, port, user, err)
		return nil, err
	}

	return db, nil
}
