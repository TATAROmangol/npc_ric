package service

import "context"

type PutRepo interface{
	PutInstitutionInfo(ctx context.Context, id int, name string, inn int) error
	PutInstitutionColumns(ctx context.Context, id int, columns []string) error
	PutMentorRequest(ctx context.Context, id int, info string) error
}

type PutService struct{
	PutRepo PutRepo
}

func (ps *PutService) PutInstitutionInfo(ctx context.Context, id int, name string, inn int) error {
	return ps.PutRepo.PutInstitutionInfo(ctx, id, name, inn)
}

func (ps *PutService) PutInstitutionColumns(ctx context.Context, id int, columns []string) error {
	return ps.PutRepo.PutInstitutionColumns(ctx, id, columns)
}

func (ps *PutService) PutMentorRequest(ctx context.Context, id int, info string) error {
	return ps.PutRepo.PutMentorRequest(ctx, id, info)
}