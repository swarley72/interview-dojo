package handler

import (
	"errors"

	corepb "github.com/swarley72/interview-dojo/proto/core"
)

var difficultyToProtoMap = map[string]corepb.Difficulty{
	"easy":   corepb.Difficulty_DIFFICULTY_EASY,
	"medium": corepb.Difficulty_DIFFICULTY_MEDIUM,
	"hard":   corepb.Difficulty_DIFFICULTY_HARD,
}

var difficultyFromProtoMap = map[corepb.Difficulty]string{
	corepb.Difficulty_DIFFICULTY_EASY:   "easy",
	corepb.Difficulty_DIFFICULTY_MEDIUM: "medium",
	corepb.Difficulty_DIFFICULTY_HARD:   "hard",
}

func difficultyToProto(d string) (corepb.Difficulty, error) {
	v, ok := difficultyToProtoMap[d]
	if !ok {
		return corepb.Difficulty_DIFFICULTY_UNSPECIFIED, errors.New("invalid difficulty")
	}

	return v, nil
}

func difficultyFromProto(d corepb.Difficulty) (string, error) {
	v, ok := difficultyFromProtoMap[d]
	if !ok {
		return "", errors.New("invalid difficulty")
	}

	return v, nil
}

var questionTypeToProtoMap = map[string]corepb.QuestionType{
	"theory":        corepb.QuestionType_QUESTION_TYPE_THEORY,
	"coding":        corepb.QuestionType_QUESTION_TYPE_CODING,
	"algorithm":     corepb.QuestionType_QUESTION_TYPE_ALGORITHM,
	"system_design": corepb.QuestionType_QUESTION_TYPE_SYSTEM_DESIGN,
}

var questionTypeFromProtoMap = map[corepb.QuestionType]string{
	corepb.QuestionType_QUESTION_TYPE_THEORY:        "theory",
	corepb.QuestionType_QUESTION_TYPE_CODING:        "coding",
	corepb.QuestionType_QUESTION_TYPE_ALGORITHM:     "algorithm",
	corepb.QuestionType_QUESTION_TYPE_SYSTEM_DESIGN: "system_design",
}

func questionTypeToProto(t string) (corepb.QuestionType, error) {
	v, ok := questionTypeToProtoMap[t]
	if !ok {
		return corepb.QuestionType_QUESTION_TYPE_UNSPECIFIED, errors.New("invalid question type")
	}

	return v, nil
}

func questionTypeFromProto(t corepb.QuestionType) (string, error) {
	v, ok := questionTypeFromProtoMap[t]
	if !ok {
		return "", errors.New("invalid question type")
	}

	return v, nil
}

func questionResponseFromProto(q *corepb.Question, p *corepb.UserProgress) (*QuestionResponse, error) {
	questionType, err := questionTypeFromProto(q.Type)
	if err != nil {
		return nil, err
	}

	difficulty, err := difficultyFromProto(q.Difficulty)
	if err != nil {
		return nil, err
	}

	var progressInfo *ProgressInfo
	if p != nil {
		progressInfo = &ProgressInfo{
			Repetitions:  p.Repetitions,
			EaseFactor:   p.EaseFactor,
			IntervalDays: p.IntervalDays,
			NextReviewAt: p.NextReviewAt.AsTime(),
		}
	}

	return &QuestionResponse{
		ID:             q.Id,
		Title:          q.Title,
		Type:           questionType,
		Difficulty:     difficulty,
		TagIDs:         q.TagIds,
		ContentMD:      q.ContentMd,
		AnswerMD:       q.AnswerMd,
		ExcalidrawJSON: q.ExcalidrawJson,
		Progress:       progressInfo,
		Verified:       q.Verified,
	}, nil
}

var answerToProtoMap = map[string]corepb.AnswerQuality{
	"again": corepb.AnswerQuality_ANSWER_QUALITY_AGAIN,
	"easy":  corepb.AnswerQuality_ANSWER_QUALITY_EASY,
	"good":  corepb.AnswerQuality_ANSWER_QUALITY_GOOD,
	"hard":  corepb.AnswerQuality_ANSWER_QUALITY_HARD,
}

var answerFromProtoMap = map[corepb.AnswerQuality]string{
	corepb.AnswerQuality_ANSWER_QUALITY_AGAIN: "again",
	corepb.AnswerQuality_ANSWER_QUALITY_EASY:  "easy",
	corepb.AnswerQuality_ANSWER_QUALITY_GOOD:  "good",
	corepb.AnswerQuality_ANSWER_QUALITY_HARD:  "hard",
}

func answerToProto(a string) (corepb.AnswerQuality, error) {
	v, ok := answerToProtoMap[a]
	if !ok {
		return corepb.AnswerQuality_ANSWER_QUALITY_AGAIN, errors.New("unknown answer")
	}
	return v, nil
}

func answerFromProto(a corepb.AnswerQuality) (string, error) {
	v, ok := answerFromProtoMap[a]
	if !ok {
		return "", errors.New("unknown answer")
	}
	return v, nil
}

func ptr[T any](v T) *T { return &v }
