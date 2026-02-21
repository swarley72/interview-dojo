package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	corepb "github.com/swarley72/interview-dojo/proto/core"
	"google.golang.org/protobuf/types/known/emptypb"
)

type listTagsResponse struct {
	Tags []tagResponse `json:"tags"`
}

type tagResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type createTagRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

func (h *Handler) ListTags(w http.ResponseWriter, req *http.Request) {
	tags, err := h.coreService.ListTags(req.Context(), &emptypb.Empty{})
	if err != nil {
		handleGRPCError(w, err, "list tags failed")
		return
	}

	res := make([]tagResponse, 0, len(tags.Tags))
	for _, t := range tags.Tags {
		res = append(res, tagResponse{ID: t.Id, Name: t.Name})
	}

	writeJSON(w, http.StatusOK, listTagsResponse{Tags: res})
}

func (h *Handler) CreateTag(w http.ResponseWriter, req *http.Request) {
	var body createTagRequest
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

	tag, err := h.coreService.CreateTag(req.Context(), &corepb.CreateTagRequest{Name: body.Name})
	if err != nil {
		handleGRPCError(w, err, "create tags failed")
		return
	}

	writeJSON(w, http.StatusCreated, tagResponse{ID: tag.Id, Name: tag.Name})
}

func (h *Handler) DeleteTag(w http.ResponseWriter, req *http.Request) {
	idParam := chi.URLParam(req, "id")
	if idParam == "" {
		writeError(w, http.StatusBadRequest, "invalid tag id")
		return
	}

	tagID, err := strconv.Atoi(idParam)
	if err != nil {
		writeError(w, http.StatusBadRequest, "tag id is not a number")
		return
	}

	_, err = h.coreService.DeleteTag(req.Context(), &corepb.DeleteTagRequest{Id: int32(tagID)})
	if err != nil {
		handleGRPCError(w, err, "delete tag failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
