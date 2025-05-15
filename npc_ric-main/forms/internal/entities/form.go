package entities

type Form struct {
	Id int `json:"id"`
	Info []string `json:"info"`
	InstitutionId int `json:"institution_id"`
}