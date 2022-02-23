package service

import (
	"context"
	"errors"

	"github.com/NTHU-LSALAB/NTHU-Distributed-System/modules/comment/dao"
	"github.com/NTHU-LSALAB/NTHU-Distributed-System/modules/comment/pb"
	"github.com/google/uuid"
)

type service struct {
	pb.UnimplementedCommentServer

	commentDAO dao.CommentDAO
}

func NewService(commentDAO dao.CommentDAO) *service {
	return &service{
		commentDAO: commentDAO,
	}
}

func (s *service) Healthz(ctx context.Context, req *pb.HealthzRequest) (*pb.HealthzResponse, error) {
	return &pb.HealthzResponse{Status: "ok"}, nil
}

func (s *service) ListComment(ctx context.Context, req *pb.ListCommentRequest) (*pb.ListCommentResponse, error) {
	comments, err := s.commentDAO.List(ctx, req.GetVideoId(), int(req.GetLimit()), int(req.GetSkip()))
	if err != nil {
		return nil, err
	}

	pbComments := make([]*pb.CommentInfo, 0, len(comments))
	for _, comment := range comments {
		pbComments = append(pbComments, comment.ToProto())
	}

	return &pb.ListCommentResponse{Comments: pbComments}, nil
}

func (s *service) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	var comment = &dao.Comment{
		VideoID: req.GetVideoId(),
		Content: req.GetContent(),
	}
	commentID, err := s.commentDAO.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &pb.CreateCommentResponse{Id: commentID.String()}, nil
}

func (s *service) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	commentID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	var comment = &dao.Comment{
		ID:      commentID,
		Content: req.GetContent(),
	}
	err = s.commentDAO.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateCommentResponse{}, nil
}

func (s *service) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	commentID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	if err := s.commentDAO.Delete(ctx, commentID); err != nil {
		if errors.Is(err, dao.ErrCommentNotFound) {
			return nil, ErrCommentNotFound
		}

		return nil, err
	}

	return &pb.DeleteCommentResponse{}, nil
}
