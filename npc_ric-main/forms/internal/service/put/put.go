package put

import "context"

//go:generate mockgen -source=put.go -destination=./mocks/put_mock.go -package=mocks

type PutRepo interface{
	PutInstitutionInfo(ctx context.Context, id int, name string, inn int) error
	PutInstitutionColumns(ctx context.Context, id int, columns []string) error
	PutMentor(ctx context.Context, id int, name string) error
	DeleteForms(ctx context.Context, institution_id int) error
}

type PutService struct{
	PutRepo PutRepo
}

func NewPutService(repo PutRepo) *PutService {
	return &PutService{
		PutRepo: repo,
	}
}

func (ps *PutService) PutInstitutionInfo(ctx context.Context, id int, name string, inn int) error {
	return ps.PutRepo.PutInstitutionInfo(ctx, id, name, inn)
}

func (ps *PutService) PutInstitutionColumns(ctx context.Context, id int, columns []string) error {
	if err := ps.PutRepo.DeleteForms(ctx, id); err != nil {
		return err
	}

	return ps.PutRepo.PutInstitutionColumns(ctx, id, columns)
}

func (ps *PutService) PutMentor(ctx context.Context, id int, name string) error {
	return ps.PutRepo.PutMentor(ctx, id, name)
}