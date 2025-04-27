package post

import "context"

//go:generate mockgen -source=post.go -destination=./mocks/post_mock.go -package=mocks

type PostRepo interface {
	PostInstitution(ctx context.Context, name string, inn int, columns []string) (int, error)
	PostMentor(ctx context.Context, name string) (int, error)
	PostForm(ctx context.Context, info []string, institutionId int, mentorId int) (int, error)
}

type PostService struct{
	PostRepo PostRepo
}

func NewPostService(postRepo PostRepo) *PostService {
	return &PostService{
		PostRepo: postRepo,
	}
}

func (ps *PostService) PostInstitution(ctx context.Context, name string, inn int, columns []string) (int, error) {
	return ps.PostRepo.PostInstitution(ctx, name, inn, columns)
}

func (ps *PostService) PostMentor(ctx context.Context, name string) (int, error) {
	return ps.PostRepo.PostMentor(ctx, name)
}	

func (ps *PostService) PostForm(ctx context.Context, info []string, institutionId int, mentorId int) (int, error) {
	return ps.PostRepo.PostForm(ctx, info, institutionId, mentorId)
}