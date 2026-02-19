package server

import (
	"context"

	corepb "github.com/swarley72/interview-dojo/proto/core"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GRPCServer) CreateTag(ctx context.Context, req *corepb.CreateTagRequest) (*corepb.Tag, error) {
	tag, err := g.tagService.CreateTag(ctx, req.Name)
	if err != nil {
		return nil, mapError(err)
	}

	return &corepb.Tag{Id: tag.ID, Name: tag.Name}, nil
}

func (g *GRPCServer) DeleteTag(ctx context.Context, req *corepb.DeleteTagRequest) (*emptypb.Empty, error) {
	err := g.tagService.DeleteTag(ctx, req.Id)
	if err != nil {
		return nil, mapError(err)
	}

	return nil, nil
}

func (g *GRPCServer) ListTags(ctx context.Context, _ *emptypb.Empty) (*corepb.ListTagsResponse, error) {
	tags, err := g.tagService.ListTags(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	res := make([]*corepb.Tag, 0, len(tags))
	for _, t := range tags {
		res = append(res, &corepb.Tag{Id: t.ID, Name: t.Name})
	}

	return &corepb.ListTagsResponse{Tags: res}, nil
}
