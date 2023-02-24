package user

import (
	"deep-ai-server/app/db"
	"deep-ai-server/app/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func IsPasswordValid(u *model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func Login(u *model.User) {
	u.LastLoginTime = time.Now()
	db.DB().Model(u).Update(&u)
}
