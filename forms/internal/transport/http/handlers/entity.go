package handlers

import "forms/internal/entities"

type GetInstitutionsResponse []entities.Institution

type GetMentorsIdResponse []entities.Mentor

type SendFormRequest struct {
	Institution string   `json:"institution"`
	Info        []string `json:"info"`
}

type SendFormResponse struct {
	ID int `json:"id"`
}