package entities

type Institution struct {
	Id int `json:"id"`
	Name string `json:"name"`
	INN int `json:"inn"`
	Columns []string 
}