package delete

import "context"

//go:generate mockgen -source=delete.go -destination=./mocks/delete_mock.go -package=mocks

type DeleteRepo interface{
	DeleteInstitution(ctx context.Context, institutionId int) error
	DeleteMentor(ctx context.Context, mentorId int) error
}

type DeleteService struct{
	DeleteRepo DeleteRepo
}

func NewDeleteService(repo DeleteRepo) *DeleteService {
	return &DeleteService{
		DeleteRepo: repo,
	}
}

func (ds *DeleteService) DeleteInstitution(ctx context.Context, institutionId int) error {
	return ds.DeleteRepo.DeleteInstitution(ctx, institutionId)
}

func (ds *DeleteService) DeleteMentor(ctx context.Context, mentorId int) error {
	return ds.DeleteRepo.DeleteMentor(ctx, mentorId)
}