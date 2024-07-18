package service

import (
	"context"
	"errors"
	"session-server/entity/errs"
	"session-server/entity/pb"
	"session-server/logic"
)

type SessionServerImpl struct {
	sessionService *logic.SessionService
	pb.UnimplementedSessionServer
}

func NewSessionServer(sessionService *logic.SessionService) pb.SessionServer {
	return &SessionServerImpl{sessionService: sessionService}
}

func (s *SessionServerImpl) Create(ctx context.Context, req *pb.CreateReq) (*pb.CreateResp, error) {
	sessionId, err := s.sessionService.Create(ctx, req.MaxInactiveInterval, req.Data, req.Attributes)
	if err != nil {
		return nil, err
	}
	return &pb.CreateResp{SessionId: sessionId}, nil
}

func (s *SessionServerImpl) SetAttribute(ctx context.Context, req *pb.SetAttributeReq) (*pb.SetAttributeResp, error) {
	err := s.sessionService.SetAttribute(ctx, req.SessionId, req.Key, req.Value)
	if err != nil {
		return nil, err
	}
	return &pb.SetAttributeResp{}, nil
}

func (s *SessionServerImpl) GetAttribute(ctx context.Context, req *pb.GetAttributeReq) (*pb.GetAttributeResp, error) {
	value, err := s.sessionService.GetAttribute(ctx, req.GetSessionId(), req.GetKey())
	if err != nil {
		if errors.Is(err, errs.SessionInvalid) {
			return &pb.GetAttributeResp{SessionInvalid: true}, nil
		}
		return nil, err
	}
	return &pb.GetAttributeResp{SessionInvalid: false, Value: value}, err
}

func (s *SessionServerImpl) Get(ctx context.Context, req *pb.GetReq) (*pb.GetResp, error) {
	data, all, err := s.sessionService.Get(ctx, req.GetSessionId())
	if err != nil {
		if errors.Is(err, errs.SessionInvalid) {
			return &pb.GetResp{SessionInvalid: true}, nil
		}
		return nil, err
	}
	return &pb.GetResp{SessionInvalid: false, Data: data, Attributes: all}, err
}
func (s *SessionServerImpl) GetData(ctx context.Context, req *pb.GetDataReq) (*pb.GetDataResp, error) {
	data, err := s.sessionService.GetData(ctx, req.GetSessionId())
	if err != nil {
		if errors.Is(err, errs.SessionInvalid) {
			return &pb.GetDataResp{SessionInvalid: true}, nil
		}
		return nil, err
	}
	return &pb.GetDataResp{SessionInvalid: false, Data: data}, err
}

func (s *SessionServerImpl) RemoveAttribute(ctx context.Context, req *pb.RemoveAttributeReq) (*pb.RemoveAttributeResp, error) {
	if err := s.sessionService.RemoveAttribute(ctx, req.SessionId, req.Key); err != nil {
		return nil, err
	}
	return &pb.RemoveAttributeResp{}, nil
}

func (s *SessionServerImpl) Invalidate(ctx context.Context, req *pb.InvalidateReq) (*pb.InvalidateResp, error) {
	if err := s.sessionService.Invalidate(ctx, req.SessionId); err != nil {
		return nil, err
	}
	return &pb.InvalidateResp{}, nil
}
