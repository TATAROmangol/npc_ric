package post

type PostInstitutionRequest struct {
	Name    string   `json:"name"`
	INN     int      `json:"inn"`
	Columns []string `json:"columns"`
}

type PostInstitutionResponse struct {
	Id int `json:"id"`
}

type PostMentorRequest struct {
	Name string `json:"name"`
}

type PostMentorResponse struct {
	Id int `json:"id"`
}

type PostFormRequest struct {
	Info          []string `json:"info"`
	InstitutionId int      `json:"institution_id"`
}

type PostFormResponse struct {
	Id int `json:"id"`
}
