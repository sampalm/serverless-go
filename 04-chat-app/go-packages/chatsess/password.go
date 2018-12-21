package chatsess

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func NewPassword(s string) string {
	salt := make([]byte, 10)
	rand.Read(salt)
	return password(s, salt)
}

func CheckPassword(p, cp string) bool {
	s := strings.Split(cp, "_")[0]
	salt, _ := hex.DecodeString(s)
	return password(p, salt) == cp
}

func password(s string, salt []byte) string {
	key, _ := scrypt.Key([]byte(s), salt, 32768, 8, 1, 32)
	return fmt.Sprintf("%x_%x", salt, key)
}
