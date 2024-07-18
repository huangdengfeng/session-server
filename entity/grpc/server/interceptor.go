package server

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"session-server/entity/errs"
	"session-server/entity/grpc/header"
	"session-server/entity/validator"
	"strconv"
	"time"
)

func CreateDefaultInterceptor() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(createContextHeaderInterceptor(), createAcclogInterceptor(), createErrInterceptor(), createValidateInterceptor())
}

// 上下文拦截器，需要放在第一个
func createContextHeaderInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctxHeader := &header.Header{
			InMetadata:  map[string][]string{},
			OutMetadata: map[string][]string{},
		}
		// 获取请求header
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctxHeader.InMetadata = md
		}

		ctx = context.WithValue(ctx, header.Key{}, ctxHeader)
		resp, err = handler(ctx, req)

		// 主动通过metadata.NewOutgoingContext 设置过响应头，则合并
		if md, exists := metadata.FromOutgoingContext(ctx); exists {
			ctxHeader.OutMetadata = metadata.Join(ctxHeader.OutMetadata, md)
		}
		if len(ctxHeader.OutMetadata) > 0 {
			if e := grpc.SendHeader(ctx, ctxHeader.OutMetadata); e != nil {
				return nil, e
			}
		}
		return resp, err
	}
}

func createAcclogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		now := time.Now()
		resp, err = handler(ctx, req)
		ip := ""
		p, ok := peer.FromContext(ctx)
		if ok {
			ip = p.Addr.String()
		}
		errCode := "0"
		outMetadata := ctx.Value(header.Key{}).(*header.Header).OutMetadata
		if len(outMetadata.Get(header.ErrCode)) > 0 {
			errCode = outMetadata.Get(header.ErrCode)[0]
		}
		log.Infof("server method %s client %s errcode %s cost %f ms", info.FullMethod, ip, errCode, float64(time.Now().Sub(now).Nanoseconds())/1e6)
		return resp, err
	}
}

// 全局错误拦截器
func createErrInterceptor() grpc.UnaryServerInterceptor {

	var setErrMetadata = func(ctx context.Context, e *errs.Error) {
		if ctxHeader, ok := ctx.Value(header.Key{}).(*header.Header); ok {
			// 使用扩展头传递错误信息
			ctxHeader.OutMetadata.Set(header.ErrType, strconv.Itoa(int(e.Type)))
			ctxHeader.OutMetadata.Set(header.ErrCode, strconv.Itoa(int(e.Code)))
			ctxHeader.OutMetadata.Set(header.ErrMsg, e.Msg)
		}
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return
		}
		log.Errorf("server execute [%s] err [%s]", info.FullMethod, err)

		if e, ok := err.(*errs.Error); ok {
			setErrMetadata(ctx, e)
			return nil, status.Errorf(codes.Internal, e.Error())
		}

		// 调用下游grpc 出错
		if s, ok := status.FromError(err); ok {
			remoteError := errs.RpcError.Newf(s.Message())
			setErrMetadata(ctx, remoteError)
			return nil, status.Errorf(codes.Internal, remoteError.Error())
		}

		unknown := errs.Unknown.Newf(err)
		setErrMetadata(ctx, unknown)
		return nil, status.Errorf(codes.Internal, unknown.Error())
	}
}

func createValidateInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if verifiable, ok := req.(validator.Verifiable); ok {
			e := verifiable.Validate()
			if e != nil {
				return nil, errs.BasArgs.Newf(e)
			}
		}
		return handler(ctx, req)
	}
}
