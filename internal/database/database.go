package database

import (
	"errors"
	"fmt"
	"fs/gofs_file_server/internal/files"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DbInstance struct {
	Db *sqlx.DB
	Lg *files.Logger
	ex bool
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

const (
	port   = 5432
	host   = "localhost"
	user   = "postgres"
	dbname = "backend"
)

const tables = `
CREATE TABLE IF NOT EXISTS Users 
(ID SERIAL PRIMARY KEY,
Username TEXT NOT NULL,
Password TEXT NOT NULL,
Email TEXT NOT NULL,
Salt TEXT NOT NULL);

`

var gdb DbInstance

func InitDb() {
	gdb.ex = false
}

func GetInstanceDb() DbInstance {
	if !gdb.ex {
		gdb.Db, gdb.Lg = connectPostgres()
		gdb.ex = true
		return gdb
	}
	return gdb
}

func connectPostgres() (*sqlx.DB, *files.Logger) {
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

	db.MustExec(tables)

	/*_, err = db.Begin()
	if err != nil {
		panic(err)
	}

	var stmt *sql.Stmt
	stmt, err = db.Prepare("INSERT INTO Users (username, password) VALUES($1, $2);")
	if err != nil {
		fmt.Println("hre1")
		panic(err)
	}
	_, err = stmt.Exec("Yaniv", "Yaniv123")
	if err != nil {
		fmt.Println("hre2")
		panic(err)
	}

	stmt.Close()*/
	return db, lg
}

// function to close  the db connection.
func (db *DbInstance) ClosePostgres() error {
	err := db.Db.Close() // close the db
	if err != nil {
		return err
	}
	err = db.Lg.DestroyLog() // close the logger
	if err != nil {
		return err
	}
	return nil
}

func (db *DbInstance) DoesUserExists(usr User) error {
	tx, err := db.Db.Begin()
	if err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("SELECT 1 FROM users WHERE EXISTS Username=$1")
	sb.WriteString(usr.Username)
	stmt, err := tx.Prepare(sb.String())
	if err != nil {
		return err
	}

	rows, err := stmt.Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int32
	for rows.Next() {
		count += 1
		if count >= 1 {
			return nil
		}
	}

	return errors.New("the user doesn't exists")
}

/*func (db *DbInstance) CheckCredentials(usr User) error {

}*/
