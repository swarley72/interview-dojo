package server

import (
	"errors"

	"github.com/swarley72/interview-dojo/core-service/internal/service"
	corepb "github.com/swarley72/interview-dojo/proto/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	corepb.UnimplementedCoreServiceServer
	questionService service.QuestionService
	tagService      service.TagService
	progressService service.UserProgressService
}

func NewGRPCServer(
	questionService service.QuestionService,
	tagService service.TagService,
	progressService service.UserProgressService,
) *GRPCServer {
	return &GRPCServer{
		questionService: questionService,
		tagService:      tagService,
		progressService: progressService,
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrQuestionNotFound):
		return status.Error(codes.NotFound, "question not found")
	case errors.Is(err, ErrUnknownEnumValue):
		return status.Error(codes.Internal, "unknown enum value")
	case errors.Is(err, service.ErrProgressNotFound):
		return status.Error(codes.NotFound, "user progress not found")
	case errors.Is(err, service.ErrNoQuestionsAvailable):
		return status.Error(codes.NotFound, "question not found")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
