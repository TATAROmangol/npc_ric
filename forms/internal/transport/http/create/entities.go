package httpcreate

import "forms/internal/entities"

type GetInstitutionsResponse struct {
	Institutions []entities.Institution `json:"institutions"`
}

