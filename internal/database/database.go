package database

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Yanivshus/gofs/internal/crypto"
	"github.com/Yanivshus/gofs/internal/files"
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

// user for internal use (including all attributes that doesn't get passed in the json)
type IUser struct {
	Id       uint64
	Username string
	Password []byte
	Email    string
	Salt     string
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

// works like a singleton.
func GetInstanceDb() DbInstance {
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

	lg := files.CreateLogger("DB", "db.log")
	go lg.KeepLogger()

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

func (db *DbInstance) DoesUserExists(usr User, where string) (bool, error) {
	tx, err := db.Db.Begin()
	if err != nil {
		return false, err
	}

	var sb strings.Builder
	sb.WriteString("SELECT EXISTS(SELECT 1 FROM ")
	sb.WriteString(where)
	sb.WriteString(" WHERE username=$1);")

	stmt, err := tx.Prepare(sb.String())
	if err != nil {
		return false, err
	}
	defer tx.Commit()

	rows, err := stmt.Query(usr.Username)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var res bool
		if err := rows.Scan(&res); err != nil {
			return true, err
		}
		return res, nil
	}

	return true, nil
}

func (db *DbInstance) GetUserData(username string) (*IUser, error) {
	tx, err := db.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	stmt, err := tx.Prepare("SELECT * FROM users WHERE username=$1;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		return nil, err
	}

	var usr IUser

	for rows.Next() {
		if err := rows.Scan(&usr.Id, &usr.Username, &usr.Password, &usr.Email, &usr.Salt); err != nil {
			return nil, err
		}
	}

	return &usr, nil
}

func (db *DbInstance) AddUserToPending(usr User) error {
	tx, err := db.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit() // commit the transaction at end of life of tx.

	var sb strings.Builder
	sb.WriteString("INSERT INTO Pending (Username, Password, Email, Salt) VALUES ($1,$2,$3,$4);")
	stmt, err := tx.Prepare(sb.String())
	if err != nil {
		return err
	}

	salt, err := crypto.GenSalt()
	if err != nil {
		return err
	}

	res, err := stmt.Exec(usr.Username, crypto.HashArgon(usr.Password, salt), usr.Email, salt)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("problem inserting pending user")
	}

	var log strings.Builder
	log.WriteString("Added user to pending, username: ")
	log.WriteString(usr.Username)
	db.Lg.LogDb(log.String())

	return nil
}

func (db *DbInstance) LoginUser(usr User) (bool, error) {
	user, err := db.GetUserData(usr.Username)
	if err != nil {
		return false, errors.New("user doesn't exists")
	}

	if bytes.Equal(user.Password, crypto.HashArgon(usr.Password, user.Salt)) {
		return true, nil
	}

	return false, nil
}
