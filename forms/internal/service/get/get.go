package get

import (
	"context"
	"forms/internal/entities"
)

//go:generate mockgen -source=get.go -destination=./mocks/pet_mock.go -package=mocks

type GetRepo interface {
	GetInstitutions(ctx context.Context) ([]entities.Institution, error)
	GetMentors(ctx context.Context) ([]entities.Mentor, error)
	GetInstitutionFromINN(ctx context.Context, inn int) (entities.Institution, error)
	GetFormColumns(ctx context.Context, id int) ([]string, error)
	GetFormRows(ctx context.Context, id int) ([][]string, error)
}

type GetService struct {
	GetRepo GetRepo
}

func NewGetService(repo GetRepo) *GetService {
	return &GetService{
		GetRepo: repo,
	}
}

func (gs *GetService) GetInstitutions(ctx context.Context) ([]entities.Institution, error) {
	return gs.GetRepo.GetInstitutions(ctx)
}

func (gs *GetService) GetMentors(ctx context.Context) ([]entities.Mentor, error) {
	return gs.GetRepo.GetMentors(ctx)
}

func (gs *GetService) GetInstitutionFromINN(ctx context.Context, inn int) (entities.Institution, error) {
	return gs.GetRepo.GetInstitutionFromINN(ctx, inn)
}

func (gs *GetService) GetFormColumns(ctx context.Context, id int) ([]string, error) {
	return gs.GetRepo.GetFormColumns(ctx, id)
}

func (gs *GetService) GetFormRows(ctx context.Context, id int) ([][]string, error) {
	return gs.GetRepo.GetFormRows(ctx, id)
}