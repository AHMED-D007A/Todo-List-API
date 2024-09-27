package user

import "net/http"

type UserHandler struct {
	storage *UserStorage
}

func NewUserHanlder(storage *UserStorage) *UserHandler {
	return &UserHandler{
		storage: storage,
	}
}

func (uh *UserHandler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {}

func (uh *UserHandler) VerifiyUser(w http.ResponseWriter, r *http.Request) {}
