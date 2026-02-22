package server

import (
	"context"

	"github.com/swarley72/interview-dojo/core-service/internal/repository"
	corepb "github.com/swarley72/interview-dojo/proto/core"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GRPCServer) GetQuestion(ctx context.Context, req *corepb.GetQuestionRequest) (*corepb.Question, error) {
	question, err := g.questionService.GetQuestionByID(ctx, req.QuestionId)
	if err != nil {
		return nil, mapError(err)
	}

	questionProto, err := questionToProto(&question)
	if err != nil {
		return nil, mapError(err)
	}

	return questionProto, nil
}

func (g *GRPCServer) DeleteQuestion(ctx context.Context, req *corepb.DeleteQuestionRequest) (*emptypb.Empty, error) {
	err := g.questionService.DeleteQuestion(ctx, req.QuestionId)
	if err != nil {
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}

func (g *GRPCServer) CreateQuestion(ctx context.Context, req *corepb.CreateQuestionRequest) (*corepb.Question, error) {
	questionType, err := questionTypeFromProto(req.Type)
	if err != nil {
		return nil, mapError(err)
	}

	questionDifficulty, err := difficultyFromProto(req.Difficulty)
	if err != nil {
		return nil, mapError(err)
	}

	question, err := g.questionService.CreateQuestion(ctx, repository.CreateQuestionParams{
		Title:      req.Title,
		ContentMD:  req.ContentMd,
		AnswerMD:   req.AnswerMd,
		Type:       questionType,
		Difficulty: questionDifficulty,
		TagIDs:     req.TagIds,
	})
	if err != nil {
		return nil, mapError(err)
	}

	questionProto, err := questionToProto(&question)
	if err != nil {
		return nil, mapError(err)
	}

	return questionProto, nil
}

func (g *GRPCServer) UpdateQuestion(ctx context.Context, req *corepb.UpdateQuestionRequest) (*corepb.Question, error) {
	var questionType *string
	if req.Type != nil {
		t, err := questionTypeFromProto(*req.Type)
		if err != nil {
			return nil, mapError(err)
		}
		questionType = &t
	}

	var questionDifficulty *string
	if req.Difficulty != nil {
		d, err := difficultyFromProto(*req.Difficulty)
		if err != nil {
			return nil, mapError(err)
		}
		questionDifficulty = &d
	}

	question, err := g.questionService.UpdateQuestion(ctx, req.QuestionId, repository.UpdateQuestionParams{
		Title:      req.Title,
		ContentMD:  req.ContentMd,
		AnswerMD:   req.AnswerMd,
		Type:       questionType,
		Difficulty: questionDifficulty,
		TagIDs:     req.TagIds,
	})
	if err != nil {
		return nil, mapError(err)
	}

	questionProto, err := questionToProto(&question)
	if err != nil {
		return nil, mapError(err)
	}

	return questionProto, nil
}

func (g *GRPCServer) ListQuestions(ctx context.Context, req *corepb.ListQuestionsRequest) (*corepb.ListQuestionsResponse, error) {
	var questionType *string
	if req.Type != nil {
		qt, err := questionTypeFromProto(*req.Type)
		if err != nil {
			return nil, mapError(err)
		}
		questionType = &qt
	}

	var questionDifficulty *string
	if req.Difficulty != nil {
		qd, err := difficultyFromProto(*req.Difficulty)
		if err != nil {
			return nil, mapError(err)
		}
		questionDifficulty = &qd
	}

	res, err := g.questionService.ListQuestions(ctx, repository.ListQuestionsFilters{
		Limit:      int(req.Limit),
		Offset:     int(req.Offset),
		Type:       questionType,
		Difficulty: questionDifficulty,
		TagIDs:     req.TagIds,
	})
	if err != nil {
		return nil, mapError(err)
	}

	questions := make([]*corepb.Question, 0, len(res.Questions))
	for _, q := range res.Questions {
		protoQuestionType, err := questionToProto(&q)
		if err != nil {
			return nil, mapError(err)
		}
		questions = append(questions, protoQuestionType)
	}

	return &corepb.ListQuestionsResponse{TotalCount: res.TotalCount, Questions: questions}, nil
}
