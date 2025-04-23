package admin

import (
	"log"

	"github.com/joho/godotenv"
)

type Admin struct{
	Login string `env:"ADMIN_LOGIN"`
	Password string `env:"ADMIN_PASSWORD"`
}

func MustLoad() Admin {
	admin := Admin{}
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return admin
}

func (a Admin) IsValid(login, password string) bool {
	if a.Login == login && a.Password == password {
		return true
	}
	return false
}