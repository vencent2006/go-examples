package rpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInitClientProxy(t *testing.T) {
	server := NewServer()
	msg := "hello, world"
	server.RegisterService(&UserServiceServer{Msg: msg})
	serverAddr := ":8081"
	go func() {
		err := server.Start("tcp", serverAddr)
		t.Log(err)
	}()
	time.Sleep(time.Second * 2)
	usClient := &UserService{}
	client, err := NewClient(serverAddr)
	require.NoError(t, err)
	err = client.InitService(usClient)
	require.NoError(t, err)
	resp, err := usClient.GetById(context.Background(), &GetByIdReq{Id: 123})
	require.NoError(t, err)
	assert.Equal(t, &GetByIdResp{Msg: msg}, resp)
}

func TestInitClientProxy2(t *testing.T) {
	server := NewServer()
	service := &UserServiceServer{}
	server.RegisterService(service)
	serverAddr := ":8081"
	go func() {
		err := server.Start("tcp", serverAddr)
		t.Log(err)
	}()
	time.Sleep(time.Second * 2)
	usClient := &UserService{}
	client, err := NewClient(serverAddr)
	require.NoError(t, err)
	err = client.InitService(usClient)
	require.NoError(t, err)

	testCases := []struct {
		name string
		mock func()

		wantErr  error
		wantResp *GetByIdResp
	}{
		{
			name: "no error",
			mock: func() {
				service.Err = nil
				service.Msg = "hello, world"
			},
			wantResp: &GetByIdResp{
				Msg: "hello, world",
			},
		},
		{
			name: "error",
			mock: func() {
				service.Msg = ""
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIdResp{},
			wantErr:  errors.New("mock error"),
		},
		{
			name: "both",
			mock: func() {
				service.Msg = "hello, world"
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIdResp{},
			wantErr:  errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			resp, er := usClient.GetById(context.Background(), &GetByIdReq{Id: 123})
			assert.Equal(t, tc.wantErr, er)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}
