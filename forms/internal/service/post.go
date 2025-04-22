package service

import "context"

type PostRepo interface {
	PostInstitution(ctx context.Context, name string, inn int, columns []string) (int, error)
	PostMentor(ctx context.Context, name string) (int, error)
	PostForm(ctx context.Context, institutionId int, info []string) (int, error)
}

type PostService struct{
	PostRepo PostRepo
}

func (ps *PostService) PostInstitution(ctx context.Context, name string, inn int, columns []string) (int, error) {
	return ps.PostRepo.PostInstitution(ctx, name, inn, columns)
}

func (ps *PostService) PostMentor(ctx context.Context, name string) (int, error) {
	return ps.PostRepo.PostMentor(ctx, name)
}	

func (ps *PostService) PostForm(ctx context.Context, institutionId int, info []string) (int, error) {
	return ps.PostRepo.PostForm(ctx, institutionId, info)
}