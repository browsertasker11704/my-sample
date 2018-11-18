package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spy16/droplet/internal/domain"
	"github.com/spy16/droplet/pkg/logger"
)

func addUsersAPI(logger logger.Logger, router *mux.Router, reg registration, ret retriever) {
	uc := &userController{
		Logger: logger,
		reg:    reg,
		ret:    ret,
	}

	router.HandleFunc("/v1/users/{name}", uc.get).Methods(http.MethodGet)
	router.HandleFunc("/v1/users/", uc.post).Methods(http.MethodPost)
}

type userController struct {
	logger.Logger
	reg registration
	ret retriever
}

func (uc *userController) get(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	user, err := uc.ret.Get(req.Context(), vars["name"])
	if err != nil {
		writeError(wr, err)
		return
	}

	writeResponse(wr, http.StatusOK, user)
}

func (uc *userController) post(wr http.ResponseWriter, req *http.Request) {
	user := domain.User{}
	if err := readRequest(req, &user); err != nil {
		uc.Warnf("failed to read user request: %s", err)
		writeResponse(wr, http.StatusBadRequest, err)
		return
	}

	registered, err := uc.reg.Register(req.Context(), user)
	if err != nil {
		uc.Warnf("failed to register user: %s", err)
		writeError(wr, err)
		return
	}

	uc.Infof("new user registered with id '%s'", registered.Name)
	writeResponse(wr, http.StatusCreated, registered)
}

type registration interface {
	Register(ctx context.Context, user domain.User) (*domain.User, error)
}

type retriever interface {
	Get(ctx context.Context, name string) (*domain.User, error)
}