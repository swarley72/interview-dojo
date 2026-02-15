package handler

import (
	"encoding/json"
	"net/http"

	userpb "github.com/swarley72/interview-dojo/proto/user"
)

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=1"`
}

type refreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) RefreshToken(w http.ResponseWriter, req *http.Request) {
	var body refreshTokenRequest
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	err = h.validate.Struct(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	response, err := h.userService.RefreshToken(
		req.Context(),
		&userpb.RefreshTokenRequest{
			RefreshToken: body.RefreshToken,
		},
	)
	if err != nil {
		handleGRPCError(w, err, "refresh failed")
		return
	}

	writeJSON(
		w,
		http.StatusOK, refreshTokenResponse{
			AccessToken:  response.AccessToken,
			RefreshToken: response.RefreshToken,
		},
	)
}
