package server

import (
	"context"

	corepb "github.com/swarley72/interview-dojo/proto/core"
)

func (g *GRPCServer) GetUserProgress(ctx context.Context, req *corepb.GetUserProgressRequest) (*corepb.UserProgress, error) {
	progress, err := g.progressService.GetProgress(ctx, req.UserId, req.QuestionId)
	if err != nil {
		return nil, mapError(err)
	}

	return progressToProto(&progress), nil
}

func (g *GRPCServer) RecordAnswer(ctx context.Context, req *corepb.RecordAnswerRequest) (*corepb.UserProgress, error) {
	answerQuality, err := answerQualityFromProto(req.AnswerQuality)
	if err != nil {
		return nil, mapError(err)
	}

	progress, err := g.progressService.RecordAnswer(ctx, req.UserId, req.QuestionId, answerQuality)
	if err != nil {
		return nil, mapError(err)
	}

	return progressToProto(&progress), nil
}

func (g *GRPCServer) GetNextQuestion(ctx context.Context, req *corepb.GetNextQuestionRequest) (*corepb.Question, error) {
	question, err := g.progressService.GetNextQuestion(ctx, req.UserId)
	if err != nil {
		return nil, mapError(err)
	}

	res, err := questionToProto(&question)
	if err != nil {
		return nil, mapError(err)
	}

	return res, nil
}
