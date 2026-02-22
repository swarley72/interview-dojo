package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/swarley72/interview-dojo/api-gateway/internal/middleware"
	corepb "github.com/swarley72/interview-dojo/proto/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecordAnswerRequest struct {
	Answer string `json:"answer" validate:"required,oneof=again easy good hard"`
}

type UserProgressResponse struct {
	Repetitions  int32   `json:"repetitions"`
	IntervalDays int32   `json:"interval_days"`
	EaseFactor   float32 `json:"ease_factor"`
}

func (h *Handler) GetNextQuestion(w http.ResponseWriter, req *http.Request) {
	userID := middleware.UserIDFromAuthClaims(req.Context())
	q := req.URL.Query()
	grpcReq := &corepb.GetNextQuestionRequest{UserId: userID}

	if t := q.Get("type"); t != "" {
		qt, err := questionTypeToProto(t)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid question type")
			return
		}
		grpcReq.Type = &qt
	}

	tagStrs := q["tag_id"]
	var tagIDs []int32
	for _, s := range tagStrs {
		id, err := strconv.Atoi(s)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid tag_id query")
			return
		}
		tagIDs = append(tagIDs, int32(id))
	}
	grpcReq.TagIds = tagIDs

	question, err := h.coreService.GetNextQuestion(req.Context(), grpcReq)
	if err != nil {
		handleGRPCError(w, err, "get next question failed")
		return
	}

	progress, err := h.coreService.GetUserProgress(req.Context(), &corepb.GetUserProgressRequest{UserId: userID, QuestionId: question.Id})
	if err != nil && status.Code(err) != codes.NotFound {
		handleGRPCError(w, err, "get user progress failed")
		return
	}

	res, err := questionResponseFromProto(question, progress)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "invalid enum in db")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) RecordAnswer(w http.ResponseWriter, req *http.Request) {
	userID := middleware.UserIDFromAuthClaims(req.Context())

	questionID := chi.URLParam(req, "question_id")
	if questionID == "" {
		writeError(w, http.StatusBadRequest, "invalid question id")
		return
	}

	var body RecordAnswerRequest
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

	userAnswer, err := answerToProto(body.Answer)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid answer")
		return
	}

	userProgress, err := h.coreService.RecordAnswer(req.Context(), &corepb.RecordAnswerRequest{
		UserId:        userID,
		QuestionId:    questionID,
		AnswerQuality: userAnswer,
	})
	if err != nil {
		handleGRPCError(w, err, "record answer failed")
		return
	}

	writeJSON(w, http.StatusOK, UserProgressResponse{
		Repetitions:  userProgress.Repetitions,
		IntervalDays: userProgress.IntervalDays,
		EaseFactor:   userProgress.EaseFactor,
	})
}

func (h *Handler) ResetProgress(w http.ResponseWriter, req *http.Request) {
	userID := middleware.UserIDFromAuthClaims(req.Context())

	_, err := h.coreService.ResetUserProgress(req.Context(), &corepb.ResetUserProgressRequest{UserId: userID})
	if err != nil {
		handleGRPCError(w, err, "reset progress failed")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
