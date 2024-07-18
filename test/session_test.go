package test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"session-server/entity/pb"
	"strconv"
	"testing"
)

func TestSession(t *testing.T) {
	ctx := context.Background()
	resp, err := client.Create(ctx, &pb.CreateReq{
		MaxInactiveInterval: 1800,
		Data:                []byte("1000"),
		Attributes:          nil,
	})
	assert.NoError(t, err)
	fmt.Println("sessionId:", resp.SessionId)
	setAttributeresp, err := client.SetAttribute(ctx, &pb.SetAttributeReq{
		SessionId: resp.SessionId,
		Key:       "uid",
		Value:     []byte(strconv.Itoa(1000)),
	})
	assert.NoError(t, err)
	assert.NotNil(t, setAttributeresp)

	setAttributeresp1, err := client.SetAttribute(ctx, &pb.SetAttributeReq{
		SessionId: resp.SessionId + "_not_exists",
		Key:       "uid",
		Value:     []byte(strconv.Itoa(1000)),
	})
	assert.Error(t, err)
	assert.Nil(t, setAttributeresp1)

	getAttributeResp, err := client.GetAttribute(ctx, &pb.GetAttributeReq{
		SessionId: resp.SessionId,
		Key:       "uid",
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte(strconv.Itoa(1000)), getAttributeResp.Value)

	getAttributeResp1, err := client.GetAttribute(ctx, &pb.GetAttributeReq{
		SessionId: resp.SessionId + "_not_exists",
		Key:       "uid",
	})
	assert.NoError(t, err)
	assert.True(t, getAttributeResp1.SessionInvalid)
	assert.Empty(t, getAttributeResp1.Value)

	getResp, err := client.Get(ctx, &pb.GetReq{
		SessionId: resp.SessionId,
	})
	assert.NoError(t, err)
	assert.False(t, getResp.SessionInvalid)
	assert.NotEmpty(t, getResp.Data)
	assert.NotEmpty(t, getResp.Attributes)
}

func BenchmarkCreate(b *testing.B) {
	attributes := map[string][]byte{
		"uid":      []byte("1000"),
		"username": []byte("dfenghuang"),
		"name":     []byte("黄登峰"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		resp, err := client.Create(ctx, &pb.CreateReq{
			MaxInactiveInterval: 1800,
			Attributes:          attributes})
		assert.NoError(b, err)
		assert.NotEmpty(b, resp.SessionId)
	}
}

func BenchmarkGetAllAttribute(b *testing.B) {
	attributes := map[string][]byte{
		"uid":      []byte("1000"),
		"username": []byte("dfenghuang"),
		"name":     []byte("黄登峰"),
	}
	ctx := context.Background()
	resp, err := client.Create(ctx, &pb.CreateReq{
		MaxInactiveInterval: 1800,
		Data:                []byte("1000"),
		Attributes:          attributes})
	assert.NoError(b, err)
	assert.NotEmpty(b, resp.SessionId)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getResp, err := client.Get(ctx, &pb.GetReq{
			SessionId: resp.SessionId,
		})
		assert.NoError(b, err)
		assert.False(b, getResp.SessionInvalid)
		assert.Equal(b, []byte("1000"), getResp.Data)
		assert.Equal(b, attributes, getResp.Attributes)
	}
}
