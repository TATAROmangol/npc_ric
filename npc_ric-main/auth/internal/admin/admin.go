package admin

import (
)

type Admin struct{
	cfg Config
}

func New(cfg Config) Admin {
	return Admin{
		cfg: cfg,
	}
}

func (a Admin) IsValid(login, password string) bool {
	if a.cfg.Login == login && a.cfg.Password == password {
		return true
	}
	return false
}