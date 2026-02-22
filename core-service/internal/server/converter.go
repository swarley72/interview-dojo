package server

import (
	"errors"

	"github.com/swarley72/interview-dojo/core-service/internal/repository"
	corepb "github.com/swarley72/interview-dojo/proto/core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrUnknownEnumValue = errors.New("unknown enum value")

var difficultyFromProtoMap = map[corepb.Difficulty]string{
	corepb.Difficulty_DIFFICULTY_EASY:   "easy",
	corepb.Difficulty_DIFFICULTY_MEDIUM: "medium",
	corepb.Difficulty_DIFFICULTY_HARD:   "hard",
}

func difficultyFromProto(d corepb.Difficulty) (string, error) {
	s, ok := difficultyFromProtoMap[d]
	if !ok {
		return "", ErrUnknownEnumValue
	}
	return s, nil
}

var difficultyToProtoMap = map[string]corepb.Difficulty{
	"easy":   corepb.Difficulty_DIFFICULTY_EASY,
	"medium": corepb.Difficulty_DIFFICULTY_MEDIUM,
	"hard":   corepb.Difficulty_DIFFICULTY_HARD,
}

func difficultyToProto(d string) (corepb.Difficulty, error) {
	s, ok := difficultyToProtoMap[d]
	if !ok {
		return corepb.Difficulty_DIFFICULTY_UNSPECIFIED, ErrUnknownEnumValue
	}
	return s, nil
}

var questionTypeFromProtoMap = map[corepb.QuestionType]string{
	corepb.QuestionType_QUESTION_TYPE_THEORY:        "theory",
	corepb.QuestionType_QUESTION_TYPE_CODING:        "coding",
	corepb.QuestionType_QUESTION_TYPE_ALGORITHM:     "algorithm",
	corepb.QuestionType_QUESTION_TYPE_SYSTEM_DESIGN: "system_design",
}

func questionTypeFromProto(t corepb.QuestionType) (string, error) {
	s, ok := questionTypeFromProtoMap[t]
	if !ok {
		return "", ErrUnknownEnumValue
	}
	return s, nil
}

var questionTypeToProtoMap = map[string]corepb.QuestionType{
	"theory":        corepb.QuestionType_QUESTION_TYPE_THEORY,
	"coding":        corepb.QuestionType_QUESTION_TYPE_CODING,
	"algorithm":     corepb.QuestionType_QUESTION_TYPE_ALGORITHM,
	"system_design": corepb.QuestionType_QUESTION_TYPE_SYSTEM_DESIGN,
}

func questionTypeToProto(t string) (corepb.QuestionType, error) {
	s, ok := questionTypeToProtoMap[t]
	if !ok {
		return corepb.QuestionType_QUESTION_TYPE_UNSPECIFIED, ErrUnknownEnumValue
	}
	return s, nil
}

func questionToProto(q *repository.Question) (*corepb.Question, error) {
	protoDifficulty, err := difficultyToProto(q.Difficulty)
	if err != nil {
		return nil, err
	}

	protoQuestionType, err := questionTypeToProto(q.Type)
	if err != nil {
		return nil, err
	}

	return &corepb.Question{
		Id:         q.ID,
		Title:      q.Title,
		ContentMd:  q.ContentMD,
		AnswerMd:   q.AnswerMD,
		Difficulty: protoDifficulty,
		Type:       protoQuestionType,
		TagIds:     q.TagIDs,
		CreatedAt:  timestamppb.New(q.CreatedAt),
		UpdatedAt:  timestamppb.New(q.UpdatedAt),
	}, nil
}

var answerQualityFromProtoMap = map[corepb.AnswerQuality]int{
	corepb.AnswerQuality_ANSWER_QUALITY_AGAIN: 0,
	corepb.AnswerQuality_ANSWER_QUALITY_HARD:  3,
	corepb.AnswerQuality_ANSWER_QUALITY_GOOD:  4,
	corepb.AnswerQuality_ANSWER_QUALITY_EASY:  5,
}

func answerQualityFromProto(a corepb.AnswerQuality) (int, error) {
	q, ok := answerQualityFromProtoMap[a]
	if !ok {
		return 0, ErrUnknownEnumValue
	}
	return q, nil
}

func progressToProto(p *repository.UserProgress) *corepb.UserProgress {
	return &corepb.UserProgress{
		Id:             p.ID,
		UserId:         p.UserID,
		QuestionId:     p.QuestionID,
		Repetitions:    p.Repetitions,
		EaseFactor:     float32(p.EaseFactor),
		IntervalDays:   p.IntervalDays,
		LastReviewedAt: timestamppb.New(p.LastReviewedAt),
		NextReviewAt:   timestamppb.New(p.NextReviewAt),
	}
}
