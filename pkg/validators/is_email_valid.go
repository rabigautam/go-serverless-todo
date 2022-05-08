package validators

import (
	"regexp"
)

func IsEmailValid(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$")
	if len(email) < 3||len(email) >254|| rxEmail.MatchString(email) {
		return true
	}
	return false
}
