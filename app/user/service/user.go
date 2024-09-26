package service

import (
	"context"
	"errors"
	"go_micro_todolist/app/user/repository/db/dao"
	"go_micro_todolist/app/user/repository/db/model"
	"go_micro_todolist/idl/pb"
	"go_micro_todolist/pkg/e"
	"gorm.io/gorm"
	"sync"
)

type UserSrv struct {
}

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

// 懒汉式的单例模式 lazy-load
func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

// 饿汉式的单例模式
//func GetUserSrvHugury() *UserSrv {
//	if UserSrvIns == nil {
//		return new(UserSrv)
//	}
//	return UserSrvIns
//}

func (u *UserSrv) UserLogin(ctx context.Context, in *pb.UserRequest, out *pb.UserDetailResponse) (err error) {
	out.Code = e.SUCCESS
	//查看有没有这个人
	user, err := dao.NewUserDao(ctx).FindUserByUserName(in.UserName)
	if err != nil {
		return
	}
	if user.ID == 0 {
		err = errors.New("用户不存在")
		out.Code = e.ERROR
		return
	}

	if !user.CheckPassword(in.Password) {
		err = errors.New("用户密码错误")
		out.Code = e.ERROR
		return
	}
	out.UserDetail = BuildUser(user)
	return
}

func BuildUser(item *model.User) *pb.UserModel {
	return &pb.UserModel{
		Id:        uint32(item.ID),
		UserName:  item.UserName,
		CreatedAt: item.CreatedAt.Unix(),
		UpdatedAt: item.UpdatedAt.Unix(),
	}
}

func (u *UserSrv) UserRegister(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) (err error) {
	if req.Password != req.PasswordConfirm {
		err = errors.New("两次密码输入不一致")
		return
	}
	resp.Code = e.SUCCESS
	_, err = dao.NewUserDao(ctx).FindUserByUserName(req.UserName)
	if err != nil {
		if err == gorm.ErrRecordNotFound { // 如果不存在就继续下去
			// ...continue
		} else {
			resp.Code = e.ERROR
			return
		}
	}
	user := &model.User{
		UserName: req.UserName,
	}
	// 加密密码
	if err = user.SetPassword(req.Password); err != nil {
		resp.Code = e.ERROR
		return
	}
	if err = dao.NewUserDao(ctx).CreateUser(user); err != nil {
		resp.Code = e.ERROR
		return
	}

	resp.UserDetail = BuildUser(user)
	return
}

//UserLogin(ctx context.Context, in *UserRequest, out *UserDetailResponse) error
//UserRegister(ctx context.Context, in *UserRequest, out *UserDetailResponse) error
