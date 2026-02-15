package handler

import (
	"github.com/go-playground/validator/v10"
	userpb "github.com/swarley72/interview-dojo/proto/user"
)

type Handler struct {
	userService userpb.UserServiceClient
	validate    *validator.Validate
}

func NewHandler(userService userpb.UserServiceClient, validate *validator.Validate) *Handler {
	return &Handler{userService: userService, validate: validate}
}
