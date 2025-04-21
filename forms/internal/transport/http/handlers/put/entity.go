package admin

type PutInstitutionInfoRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	INN  int    `json:"inn"`
}

type PutInstitutionColumnsRequest struct {
	Id      int      `json:"id"`
	Columns []string `json:"columns"`
}

type PutMentorRequest struct {
	Id   int    `json:"id"`
	Info string `json:"info"`
}
