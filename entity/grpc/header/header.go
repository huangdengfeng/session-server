package header

import "google.golang.org/grpc/metadata"

const (
	ErrType = "x-err-type"
	ErrCode = "x-err-code"
	ErrMsg  = "x-err-msg"
)

// Key 上下文key
type Key struct {
}

// Header 请求头，可以用于服务端和客户端
// 服务端使用：用于存放服务端接收客户端的header及即将发送给客户端的header
// 客户端使用：用于存放即将发送的header及即将接收的服务端header
type Header struct {
	// InMetadata 客户端请求携带的header
	InMetadata metadata.MD
	// OutMetadata 服务端响应的header
	OutMetadata metadata.MD
}
