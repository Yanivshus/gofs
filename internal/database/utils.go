package database

import (
	"regexp"
)

func CheckValidEmail(email string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+@[a-zA-Z]+\\.[a-zA-Z]{2,}$", email)
	return match
}
