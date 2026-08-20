package db

import (
	"database/sql"
	"fmt"
	"log"

	"dash/ent/generated"
	"dash/utils"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var Client *generated.Client

func Connect() {
	db, err := sql.Open(
		"pgx",
		fmt.Sprintf(
			"postgres://%s:%s@localhost:5432/dash?sslmode=disable",
			utils.Env.DBUser,
			utils.Env.DBPass,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	Client = generated.NewClient(
		generated.Driver(entsql.OpenDB(dialect.Postgres, db)),
	)
}
