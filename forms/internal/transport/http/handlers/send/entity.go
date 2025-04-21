package send

type SendFormRequest struct {
	InstitutionId int      `json:"institution_id"`
	Info          []string `json:"info"`
	MentorId      int      `json:"mentor_id"`
}

type SendFormResponse struct {
	Id int `json:"id"`
}