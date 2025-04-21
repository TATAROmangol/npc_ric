package get

import "forms/internal/entities"

type GetInstitutionsResponse []entities.Institution

type GetMentorsIdResponse []entities.Mentor

type GetInstitutionFromINNRequest struct {
	Inn int `json:"inn"`
}

type GetInstitutionFromINNResponse entities.Institution

type GetFormColumnsRequest struct {
	InstitutionId int `json:"institution_id"`
}

type GetFormColumnsResponse []string