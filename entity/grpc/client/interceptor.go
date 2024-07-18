package server

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"session-server/entity/errs"
	"session-server/entity/grpc/header"
	"strconv"
	"time"
)

const DefaultTimeout = 3 * time.Second

func CreateDefaultInterceptor() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(createContextHeaderInterceptor(), createAcclogInterceptor(), createDefaultTimeoutInterceptor(), createErrInterceptor())
}

func createContextHeaderInterceptor() grpc.UnaryClientInterceptor {

	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctxHeader := &header.Header{
			InMetadata:  map[string][]string{},
			OutMetadata: map[string][]string{},
		}
		// 即将发送的header
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			ctxHeader.OutMetadata = md
		}
		ctx = context.WithValue(ctx, header.Key{}, ctxHeader)
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

func createAcclogInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		now := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		errCode := errs.Success.Code
		if err != nil {
			if e, ok := err.(*errs.Error); ok {
				errCode = e.Code
			} else {
				errCode = errs.Unknown.Code
			}
		}
		log.Infof("call server %s method %s errcode %d cost %f ms", cc.Target(), method, errCode, float64(time.Now().Sub(now).Nanoseconds())/1e6)
		return err
	}
}

// createDefaultTimeoutInterceptor
// use grpc.WithUnaryInterceptor
func createDefaultTimeoutInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 检查 context 是否已经有超时设置
		if _, ok := ctx.Deadline(); !ok {
			// 没有设置超时，添加默认超时
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, DefaultTimeout)
			defer cancel()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func createErrInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 创建用于接收 headerMD
		headerMD := metadata.MD{}
		// 添加 Header的CallOption
		opts = append(opts, grpc.Header(&headerMD))
		// 调用实际的 RPC 方法
		err := invoker(ctx, method, req, reply, cc, opts...)
		if ctxHeader, ok := ctx.Value(header.Key{}).(*header.Header); ok {
			ctxHeader.InMetadata = headerMD
		}

		if err == nil {
			return nil
		}
		log.Errorf("client call [%s] err [%s]", method, err)
		code := headerMD.Get(header.ErrCode)
		etype := headerMD.Get(header.ErrType)
		msg := headerMD.Get(header.ErrMsg)
		hasErrorCode := len(code) > 0 && len(etype) > 0 && len(msg) > 0
		if hasErrorCode {
			e := &errs.Error{}
			if i, err := strconv.Atoi(etype[0]); err != nil {
				log.Errorf("client convert errtype [%s] , error [%s]", etype[0], err)
				return errs.RpcError.Newf(err)
			} else {
				e.Type = int32(i)
			}

			if i, err := strconv.Atoi(code[0]); err != nil {
				log.Errorf("client convert errcode [%s] , error [%s]", code[0], err)
				return errs.RpcError.Newf(err)
			} else {
				e.Code = int32(i)
			}
			e.Msg = msg[0]
			return e
		}
		return errs.RpcError.Newf(err)
	}
}
