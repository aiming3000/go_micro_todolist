package rpc

import (
	"context"
	"errors"

	"go_micro_todolist/idl/pb"
	"go_micro_todolist/pkg/e"
)

// UserLogin 用户登陆
func UserLogin(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDetailResponse, err error) {
	resp, err = UserService.UserLogin(ctx, req)
	if err != nil {
		return
	}

	if resp.Code != e.SUCCESS {
		err = errors.New(e.GetMsg(int(resp.Code)))
		return
	}

	return
}

// UserRegister 用户注册
func UserRegister(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDetailResponse, err error) {
	resp, err = UserService.UserRegister(ctx, req)
	if err != nil {
		return
	}

	return
}