package service

type Repo interface {
	GetRepo
	PutRepo
	PostRepo
	DeleteRepo
}

type Services struct {
	GetService
	PutService
	PostService
	DeleteService
}

func NewServices(repo Repo) *Services {
	return &Services{
		GetService: GetService{GetRepo: repo},
		PutService: PutService{PutRepo: repo},
		PostService: PostService{PostRepo: repo},
		DeleteService: DeleteService{DeleteRepo: repo},
	}
}