package main

import (
	"github.com/GUO-xjtu/pb_gen/gen/api/user"
	"github.com/gin-gonic/gin"
	"github.com/xiangqin/user_api/common"
	"github.com/xiangqin/user_api/handlers"
)

const (
	Prefix                  = "/xiangqin/user/"
)


func xiangQinUser(r *gin.Engine) {
	r.Any(Prefix+"register", common.API(
		handlers.UserRegister,
		&user.UserRegisterRequest{},
		"xiangqin",
		"user_register",
		))
}
