package database

import (
	"fmt"
	"fs/gofs_file_server/internal/files"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DbInstance struct {
	Db     *sqlx.DB
	Lg *files.Logger
}

const (
	port   = 5432
	host   = "localhost"
	user   = "postgres"
	dbname = "backend"
)

func ConnectPostgres() DbInstance {
	err := godotenv.Load(".env")
	if err != nil {
		panic("couldn't load connection details")
	}

	password := os.Getenv("DB_PASS")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("DB Connection started successfully")

	lg := files.Create_logger("DB", "db.log")
	go lg.Keep_logger()

	return DbInstance{
		db,
		lg,
	}
}
