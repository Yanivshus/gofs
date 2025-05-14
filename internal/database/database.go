package database

import (
	"errors"
	"fmt"
	"fs/gofs_file_server/internal/crypto"
	"fs/gofs_file_server/internal/files"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
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
Password BYTEA NOT NULL,
Email TEXT NOT NULL,
Salt TEXT NOT NULL);


CREATE TABLE IF NOT EXISTS Pending 
(ID SERIAL PRIMARY KEY,
Username TEXT NOT NULL,
Password BYTEA NOT NULL,
Email TEXT NOT NULL,
Salt TEXT NOT NULL);
`

var gdb DbInstance

func InitDb() {
	gdb.ex = false
}

func GetInstanceDb() DbInstance {
	fmt.Println("fdf")
	if !gdb.ex {
		gdb.Db, gdb.Lg = connectPostgres()
		gdb.ex = true
		return gdb
	}
	return gdb
}

func connectPostgres() (*sqlx.DB, *files.Logger) {
	/*err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}*/

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
	sb.WriteString("SELECT 1 FROM Pending WHERE EXISTS Username=$1")
	stmt, err := tx.Prepare(sb.String())
	if err != nil {
		return err
	}

	rows, err := stmt.Query(usr.Username)
	if err != nil {
		return err
	}
	defer rows.Close()
	defer tx.Commit()

	var count int32
	for rows.Next() {
		count += 1
		if count >= 1 {
			return errors.New("the user exists")
		}
	}

	return nil
}

func (db *DbInstance) AddUserToPending(usr User) error {
	tx, err := db.Db.Begin()
	if err != nil {
		fmt.Println("here1")
		return err
	}
	defer tx.Commit() // commit the transaction at end of life of tx.

	var sb strings.Builder
	sb.WriteString("INSERT INTO Pending (Username, Password, Email, Salt) VALUES ($1,$2,$3,$4)")
	stmt, err := tx.Prepare(sb.String())
	if err != nil {
		fmt.Println("here2")
		return err
	}

	salt, err := crypto.GenSalt()
	if err != nil {
		fmt.Println("here3")
		return err
	}

	res, err := stmt.Exec(usr.Username, crypto.HashArgon(usr.Password, salt), usr.Email, salt)
	if err != nil {
		fmt.Println("here4")
		fmt.Println(err.Error())
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		fmt.Println("here5")
		return errors.New("problem inserting pending user")
	}

	return nil
}

/*func (db *DbInstance) CheckCredentials(usr User) error {

}*/
