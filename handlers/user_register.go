package handlers

import (
	"github.com/GUO-xjtu/pb_gen/gen/api/user"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang/protobuf/proto"
	"github.com/xiangqin/user_api/common"
	userRPC "github.com/xiangqin/user_api/kitex_gen/user"
)


func UserRegister(ctx *common.RequestContext) (proto.Message, error) {
	params := ctx.Params().(*user.UserRegisterRequest)

	req := userRPC.RegisterRequest{
		Name:     params.UserName,
		Gender:   params.Gender,
		PhoneNum: params.PhoneNum,
	}
	resp, err := common.UserCli.UserRegister(ctx.GetCtx(), &req)
	if err != nil {
		logger.CtxErrorf(ctx.GetCtx(), "UserRegister failed, err:%v", err)
		return &user.UserRegisterResponse{
			Data: nil,
		}, err
	}
	logger.CtxInfof(ctx.GetCtx(), "UserRegister success, req:%+v", req)

	return &user.UserRegisterResponse{
		Data: &user.UserRegisterResponse_ResponseData{
			UserInfo: &user.UserRegisterResponse_ResponseData_UserInfo{
				UserId: resp.GetUserID(),
			},
		},
	}, nil
}
