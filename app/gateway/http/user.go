package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_micro_todolist/app/gateway/rpc"
	"go_micro_todolist/idl/pb"
	"go_micro_todolist/pkg/ctl"
	log "go_micro_todolist/pkg/logger"
	"go_micro_todolist/pkg/utils"
	"go_micro_todolist/types"
	"net/http"
)

// UserRegisterHandler 用户注册
func UserRegisterHandler(ctx *gin.Context) {
	fmt.Println(111)
	var req pb.UserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "UserRegister Bind 绑定参数失败"))
	}
	//fmt.Println(req)
	userResp, err := rpc.UserRegister(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("UserRegister:%v", err)
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserRegister RPC 调用失败"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, userResp))

}

// UserLoginHandler 用户登录
func UserLoginHandler(ctx *gin.Context) {
	var req pb.UserRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "UserLogin Bind 绑定参数失败"))
		return
	}
	userResp, err := rpc.UserLogin(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserLogin RPC 调用失败"))
		return
	}
	token, err := utils.GenerateToken(uint(userResp.UserDetail.Id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "GenerateToken 失败"))
		return
	}
	res := &types.TokenData{
		User:  userResp,
		Token: token,
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, res))
}
