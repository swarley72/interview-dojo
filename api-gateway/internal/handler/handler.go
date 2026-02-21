package handler

import (
	"github.com/go-playground/validator/v10"
	corepb "github.com/swarley72/interview-dojo/proto/core"
	userpb "github.com/swarley72/interview-dojo/proto/user"
)

type Handler struct {
	userService userpb.UserServiceClient
	coreService corepb.CoreServiceClient
	validate    *validator.Validate
}

func NewHandler(
	userService userpb.UserServiceClient,
	coreService corepb.CoreServiceClient,
	validate *validator.Validate,
) *Handler {
	return &Handler{userService: userService, coreService: coreService, validate: validate}
}
