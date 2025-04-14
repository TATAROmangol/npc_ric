package handlers

type Service interface {
	Sender
	Creater
	Getter
}

type Handler struct{
	srv Service
}

func NewHandler(srv Service) *Handler {
	return &Handler{
		srv: srv,
	}
}