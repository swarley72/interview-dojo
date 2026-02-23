package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/swarley72/interview-dojo/api-gateway/internal/middleware"
	corepb "github.com/swarley72/interview-dojo/proto/core"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateQuestionRequest struct {
	Title          string  `json:"title" validate:"required,min=1"`
	ContentMD      *string `json:"content_md"`
	AnswerMD       *string `json:"answer_md"`
	ExcalidrawJSON *string `json:"excalidraw_json"`
	Difficulty     string  `json:"difficulty" validate:"required,oneof=easy medium hard"`
	Type           string  `json:"type" validate:"required,oneof=theory coding algorithm system_design"`
	TagIDs         []int32 `json:"tag_ids"`
	Verified       bool    `json:"verified"`
}

type UpdateQuestionRequest struct {
	Title          *string `json:"title"`
	ContentMD      *string `json:"content_md"`
	AnswerMD       *string `json:"answer_md"`
	ExcalidrawJSON *string `json:"excalidraw_json"`
	Difficulty     *string `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	Type           *string `json:"type" validate:"omitempty,oneof=theory coding algorithm system_design"`
	Verified       *bool   `json:"verified"`
	TagIDs         []int32 `json:"tag_ids"`
}

type QuestionResponse struct {
	ID             string        `json:"id"`
	Title          string        `json:"title"`
	ContentMD      *string       `json:"content_md"`
	AnswerMD       *string       `json:"answer_md"`
	Difficulty     string        `json:"difficulty"`
	Type           string        `json:"type"`
	ExcalidrawJSON *string       `json:"excalidraw_json"`
	TagIDs         []int32       `json:"tag_ids"`
	Progress       *ProgressInfo `json:"progress"`
	Verified       bool          `json:"verified"`
}

type ListQuestionsQueryParams struct {
	Difficulty string `validate:"omitempty,oneof=easy medium hard"`
	Type       string `validate:"omitempty,oneof=theory coding algorithm system_design"`
	Page       int32  `validate:"omitempty,min=1"`
	Limit      int32  `validate:"omitempty,min=1,max=100"`
	TagIDs     []int32
	Query      *string
	Verified   *bool
}

type QuestionResponseShort struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Difficulty string  `json:"difficulty"`
	Type       string  `json:"type"`
	TagIDs     []int32 `json:"tag_ids"`
	Verified   bool    `json:"verified"`
}

type ProgressInfo struct {
	Repetitions  int32     `json:"repetitions"`
	EaseFactor   float32   `json:"ease_factor"`
	IntervalDays int32     `json:"interval_days"`
	NextReviewAt time.Time `json:"next_review_at"`
}

type ListQuestionsResponse struct {
	TotalCount int32                    `json:"total_count"`
	Items      []*QuestionResponseShort `json:"items"`
}

func (h *Handler) CreateQuestion(w http.ResponseWriter, req *http.Request) {
	var body CreateQuestionRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	difficultyProto, err := difficultyToProto(body.Difficulty)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid difficulty")
		return
	}

	typeProto, err := questionTypeToProto(body.Type)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid question type")
		return
	}

	question, err := h.coreService.CreateQuestion(req.Context(), &corepb.CreateQuestionRequest{
		Title:          body.Title,
		ContentMd:      body.ContentMD,
		AnswerMd:       body.AnswerMD,
		ExcalidrawJson: body.ExcalidrawJSON,
		Difficulty:     difficultyProto,
		Type:           typeProto,
		TagIds:         body.TagIDs,
		Verified:       body.Verified,
	})
	if err != nil {
		handleGRPCError(w, err, "failed create question")
		return
	}

	res, err := questionResponseFromProto(question, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "invalid enum in db")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) UpdateQuestion(w http.ResponseWriter, req *http.Request) {
	questionID := chi.URLParam(req, "id")
	if questionID == "" {
		writeError(w, http.StatusBadRequest, "invalid question id")
		return
	}

	var body UpdateQuestionRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var newQuestionType *corepb.QuestionType
	if body.Type != nil {
		t, err := questionTypeToProto(*body.Type)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid question type")
			return
		}
		newQuestionType = &t
	}

	var newQuestionDifficulty *corepb.Difficulty
	if body.Difficulty != nil {
		d, err := difficultyToProto(*body.Difficulty)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid question difficulty")
			return
		}
		newQuestionDifficulty = &d
	}

	question, err := h.coreService.UpdateQuestion(req.Context(), &corepb.UpdateQuestionRequest{
		QuestionId:     questionID,
		Title:          body.Title,
		Type:           newQuestionType,
		Difficulty:     newQuestionDifficulty,
		AnswerMd:       body.AnswerMD,
		ContentMd:      body.ContentMD,
		ExcalidrawJson: body.ExcalidrawJSON,
		TagIds:         body.TagIDs,
		Verified:       body.Verified,
	})
	if err != nil {
		handleGRPCError(w, err, "update question failed")
		return
	}

	res, err := questionResponseFromProto(question, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "invalid enum in db")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) DeleteQuestion(w http.ResponseWriter, req *http.Request) {
	idParam := chi.URLParam(req, "id")
	if idParam == "" {
		writeError(w, http.StatusBadRequest, "invalid question id")
		return
	}

	_, err := h.coreService.DeleteQuestion(req.Context(), &corepb.DeleteQuestionRequest{QuestionId: idParam})
	if err != nil {
		handleGRPCError(w, err, "delete question failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetQuestion(w http.ResponseWriter, req *http.Request) {
	questionID := chi.URLParam(req, "id")
	if questionID == "" {
		writeError(w, http.StatusBadRequest, "invalid question id")
		return
	}
	userID := middleware.UserIDFromAuthClaims(req.Context())
	g, ctx := errgroup.WithContext(req.Context())

	var question *corepb.Question
	var progress *corepb.UserProgress

	g.Go(func() error {
		var err error
		question, err = h.coreService.GetQuestion(ctx, &corepb.GetQuestionRequest{QuestionId: questionID})
		return err
	})

	g.Go(func() error {
		var err error
		progress, err = h.coreService.GetUserProgress(ctx, &corepb.GetUserProgressRequest{QuestionId: questionID, UserId: userID})
		if status.Code(err) == codes.NotFound {
			return nil
		}
		return err
	})

	if err := g.Wait(); err != nil {
		handleGRPCError(w, err, "get question failed")
		return
	}

	res, err := questionResponseFromProto(question, progress)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "invalid enum in db")
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) ListQuestions(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
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
	verifiedStr := q.Get("verified")
	var verified *bool
	if verifiedStr != "" {
		switch verifiedStr {
		case "true":
			verified = ptr(true)
		case "false":
			verified = ptr(false)
		default:
			writeError(w, http.StatusBadRequest, "invalid verified param, expected true or false")
			return
		}
	}

	var querySearch *string
	if qs := q.Get("q"); qs != "" {
		querySearch = &qs
	}

	params := ListQuestionsQueryParams{
		Limit:      int32(limit),
		Page:       int32(page),
		Type:       q.Get("type"),
		Difficulty: q.Get("difficulty"),
		TagIDs:     tagIDs,
		Verified:   verified,
		Query:      querySearch,
	}

	if err := h.validate.Struct(params); err != nil {
		writeError(w, http.StatusBadRequest, "invalid query body")
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 20
	}

	var difficulty *corepb.Difficulty
	if params.Difficulty != "" {
		d, err := difficultyToProto(params.Difficulty)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid difficulty")
			return
		}
		difficulty = &d
	}

	var questionType *corepb.QuestionType
	if params.Type != "" {
		t, err := questionTypeToProto(params.Type)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid question type")
			return
		}
		questionType = &t
	}

	listQuestions, err := h.coreService.ListQuestions(req.Context(), &corepb.ListQuestionsRequest{
		Limit:      int32(params.Limit),
		Offset:     int32((params.Page - 1) * params.Limit),
		Difficulty: difficulty,
		Type:       questionType,
		TagIds:     params.TagIDs,
		Verified:   params.Verified,
		Query:      params.Query,
	})
	if err != nil {
		handleGRPCError(w, err, "list questions failed")
		return
	}

	var questions = make([]*QuestionResponseShort, len(listQuestions.Questions))
	for i, qst := range listQuestions.Questions {
		questionTypeString, err := questionTypeFromProto(qst.Type)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "invalid question type in db")
			return
		}

		questionDifficultyString, err := difficultyFromProto(qst.Difficulty)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "invalid difficulty in db")
			return
		}

		questions[i] = &QuestionResponseShort{
			ID:         qst.Id,
			Title:      qst.Title,
			Type:       questionTypeString,
			Difficulty: questionDifficultyString,
			TagIDs:     qst.TagIds,
			Verified:   qst.Verified,
		}
	}

	writeJSON(w, http.StatusOK, &ListQuestionsResponse{
		TotalCount: listQuestions.TotalCount,
		Items:      questions,
	})
}
