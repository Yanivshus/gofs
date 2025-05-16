package database

import (
	"regexp"

	"github.com/joho/godotenv"
)

func CheckValidEmail(email string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+@[a-zA-Z]+\\.[a-zA-Z]{2,}$", email)
	return match
}

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	return nil
}
