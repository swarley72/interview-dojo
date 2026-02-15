package handler

import (
	"net/http"

	"github.com/swarley72/interview-dojo/api-gateway/internal/middleware"
	userpb "github.com/swarley72/interview-dojo/proto/user"
)

type getProfileResponse struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	IsSuperUser bool   `json:"is_super_user"`
}

func (h *Handler) GetProfile(w http.ResponseWriter, req *http.Request) {
	userID, ok := req.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response, err := h.userService.GetUser(req.Context(), &userpb.GetUserRequest{UserId: userID})
	if err != nil {
		handleGRPCError(w, err, "failed get profile")
		return
	}

	writeJSON(w, http.StatusOK, getProfileResponse{ID: response.Id, Login: response.Login, IsSuperUser: response.IsSuperUser})
}
