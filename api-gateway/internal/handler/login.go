package handler

import (
	"encoding/json"
	"net/http"

	userpb "github.com/swarley72/interview-dojo/proto/user"
)

type loginRequest struct {
	Login    string `json:"login" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginResponse struct {
	UserID       string `json:"user_id"`
	Login        string `json:"login"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) {
	var body loginRequest
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.validate.Struct(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.userService.Login(
		req.Context(),
		&userpb.LoginRequest{
			Login:    body.Login,
			Password: body.Password,
		},
	)
	if err != nil {
		handleGRPCError(w, err, "login failed")
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		loginResponse{
			Login:        response.User.Login,
			UserID:       response.User.Id,
			AccessToken:  response.Token.AccessToken,
			RefreshToken: response.Token.RefreshToken,
		},
	)
}
