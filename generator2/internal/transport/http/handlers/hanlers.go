package handlers

type Srv interface {
	DeleteTemplate(id int) error
	GetAllTemplates() error
	UploadTemplate(id int) error
}

type Handlers struct{
	srv Srv
}

func New(srv Srv) *Handlers {
	return &Handlers{
		srv: srv,
	}
}
