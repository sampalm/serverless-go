package chatsess

import (
	"golang.org/x/crypto/bcrypt"
)

func NewPasswordBcrypt(s string) string {
	bs, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(bs)
}

func CheckPasswordBcrypt(p, cp string) bool {
	if bcrypt.CompareHashAndPassword([]byte(cp), []byte(p)) != nil {
		return false
	}
	return true
}
